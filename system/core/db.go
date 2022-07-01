package core

import (
	"fmt"
	"log"
	"reflect"
	"strings"
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
		Logger:               getDbLogger(),
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

func findEntityField(val_obj reflect.Value, type_obj reflect.Type, find_field_addr string) *reflect.StructField {
	//根据地址遍历对象属性找字段
	for i := 0; i < type_obj.Elem().NumField(); i++ {
		fieldType := type_obj.Elem().Field(i)
		filedVal := val_obj.Elem().Field(i)
		if fieldType.Type.Kind() == reflect.Struct {
			//return nil //暂不支持嵌套。。
		} else if fmt.Sprintf("%p", filedVal.Addr().Interface()) == find_field_addr { //传入的字段值地址==对象的属性字段值地址
			return &fieldType
		}
	}
	return nil

}

//获取表字段的实际名称
//username_db_field_name := core.GetFieldAlias(&adminuser, &adminuser.Username)
func GetFieldAlias(entity_struct interface{}, field interface{}) string {
	field_addr := fmt.Sprintf("%p", field)

	val_obj := reflect.ValueOf(entity_struct)
	type_obj := reflect.TypeOf(entity_struct)
	if type_obj.Kind() != reflect.Ptr || type_obj.Elem().Kind() != reflect.Struct {
		panic("第一个参数应该为数据实体entity的结构体指针，如：&user_entity")
	}
	if reflect.TypeOf(field).Kind() != reflect.Ptr {
		panic("第二个参数应该为数据实体entity的属性的指针，如：如：&user_entity.name")
	}

	//去对象里查找地址一样的属性字段
	fieldStruct := findEntityField(val_obj, type_obj, field_addr)
	if fieldStruct == nil {
		panic("entity字段不存在,不支持结构体嵌套")
	}

	err := Db.Statement.Parse(entity_struct)
	if err != nil {
		panic(err)
	}
	db_schema := (*Db.Statement).Schema
	db_field := db_schema.LookUpField(fieldStruct.Name)
	if db_field == nil {
		fmt.Printf("db_schema.Fields: %v\n", db_schema.Fields)
		panic("字段不存在,不支持结构体嵌套，" + fieldStruct.Name)
	}
	return db_field.DBName
}
func getConfigFromGromTag(tags_str, field_key string) (string, bool) {
	//Username string `gorm:"column:username;xx:yy"`
	tags_arr := strings.Split(tags_str, ";")
	for _, tag_str := range tags_arr {
		tag_arr := strings.Split(tag_str, ":")
		if tag_arr[0] == field_key {
			if len(tag_arr) != 2 {
				panic("gorm的entity配置TAG错误," + tag_str)
			}
			return tag_arr[1], true
		}
	}
	return "", false
}

//日志参数
func getDbLogger() logger.Interface {
	io := GetLogIoWriter(GetConfigString("logger.db_path")) //输出到日志文件
	db_slow_threshold := (GetConfigFloat64("logger.db_slow_threshold") * float64(time.Second))
	return logger.New(
		log.New(io, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Duration(db_slow_threshold),                     // 慢 SQL 阈值
			LogLevel:                  getDbLoggerLevel(GetConfigString("logger.db_level")), // 日志级别
			IgnoreRecordNotFoundError: true,                                                 // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                                // 禁用彩色打印
		},
	)
}

/**
配置中的日志基本和orm的日志级别对应
*/
func getDbLoggerLevel(str_level string) logger.LogLevel {
	var level = logger.Error
	logger_level, _ := zapcore.ParseLevel(str_level)
	if logger_level <= zapcore.InfoLevel {
		level = logger.Info
	} else if logger_level <= zapcore.WarnLevel {
		level = logger.Warn
	}
	return level
}
