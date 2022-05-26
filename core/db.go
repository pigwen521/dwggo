package core

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func init() {
	host := GetConfigString("database.host")
	port := GetConfigString("database.port")
	database := GetConfigString("database.database")
	username := GetConfigString("database.username")
	password := GetConfigString("database.password")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	var err error
	tmp := mysql.Open(dsn)

	Db, err = gorm.Open(tmp, &gorm.Config{
		DisableAutomaticPing: true, //默认检查数据版本，没必要再PING
		Logger:               getLogger(),
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			//NoLowerCase:   true,                              // skip the snake_casing of names
			//NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
	})
	if err != nil {
		LogError(err.Error())
		panic("orm db connect error")
	} else {
		LogInfo("init db and orm")
	}
	setConnPool()
}

//配置连接池
func setConnPool() {
	sqlDb, err := Db.DB()
	if err != nil {
		LogError(err.Error())
		panic("db.db connect error")
	}

	sqlDb.SetMaxOpenConns(GetConfigInt("database.pool_max_conn"))
	sqlDb.SetMaxIdleConns(GetConfigInt("database.pool_max_idle"))
	sqlDb.SetConnMaxIdleTime(time.Second * GetConfigDuration("database.pool_idle_time"))
	sqlDb.SetConnMaxLifetime(time.Second * GetConfigDuration("database.pool_life_time")) //必须同时设置Lifetime，IdleTime才生效
	LogInfo("init db conn pool")
}

//日志参数
func getLogger() logger.Interface {
	io := GetLogIoWriter(GetConfigString("logger.db_path")) //输出到日志文件
	db_slow_threshold := (GetConfigFloat64("logger.db_slow_threshold") * float64(time.Second))
	return logger.New(
		log.New(io, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Duration(db_slow_threshold),                   // 慢 SQL 阈值
			LogLevel:                  getLoggerLevel(GetConfigString("logger.db_level")), // 日志级别
			IgnoreRecordNotFoundError: true,                                               // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                              // 禁用彩色打印
		},
	)
}

/**
配置中的日志基本和orm的日志级别对应
*/
func getLoggerLevel(str_level string) logger.LogLevel {
	var level = logger.Error
	logger_level, _ := zapcore.ParseLevel(str_level)
	if logger_level <= zapcore.InfoLevel {
		level = logger.Info
	} else if logger_level <= zapcore.WarnLevel {
		level = logger.Warn
	}
	return level
}
