package core

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger_customs map[string]*zap.Logger

/**
自定义日志
独立日志文件log_name
*/
func LogErrorCustom(str string, log_name string) {
	loadCustomLogger(log_name).Error(str)
	if IsEnvDev() {
		fmt.Println(str)
	}
	SendWarning(str)
}
func LogInfoCustom(str string, log_name string) {
	_, file, line, _ := runtime.Caller(1)
	loadCustomLogger(log_name).Info(strings.TrimPrefix(file, APP_ROOT) + ":" + strconv.Itoa(line) + "	" + str)
	if IsEnvDev() {
		fmt.Println(str)
	}
}

//写自定义文件日志，并打印
func LogInfoCustomFmt(str string, log_name string) {
	_, file, line, _ := runtime.Caller(1)
	loadCustomLogger(log_name).Info(strings.TrimPrefix(file, APP_ROOT) + ":" + strconv.Itoa(line) + "	" + str)
	fmt.Println(str)
}

/**
记录错误日志
*/
func LogError(str string) {
	loadErrorLogger().Error(str)
	if IsEnvDev() {
		fmt.Println(str)
	}
	SendWarning(str)
}
func LogWarning(str string) {
	loadErrorLogger().Warn(str)
	if IsEnvDev() {
		fmt.Println(str)
	}
	SendWarning(str)
}
func LogErrorAndPanic(str string) {
	loadErrorLogger().Error(str)
	if IsEnvDev() {
		fmt.Println(str)
	}
	SendWarning(str)
	panic(str)
}

/**
记录info warning日志
*/
func LogInfo(str string) {
	_, file, line, _ := runtime.Caller(1)
	loadInfoLogger().Info(file + "	" + strconv.Itoa(line) + "	" + str)
	if IsEnvDev() {
		fmt.Println(str)
	}
}
func LogDebug(str string) {
	_, file, line, _ := runtime.Caller(1)
	loadInfoLogger().Debug(file + "	" + strconv.Itoa(line) + "	" + str)
	if IsEnvDev() {
		fmt.Println(str)
	}
}

/**
logger := core.LoadLogger()
logger.Info("info msg ..aa", zap.String("myurl", "my..url.."))
logger.Error("error msg ..aa", zap.String("myurl", "my..url.."))
*/
func loadInfoLogger() *zap.Logger {
	return loadLogger(GetConfigString("logger.info_path"))
}
func loadErrorLogger() *zap.Logger {
	return loadLogger(GetConfigString("logger.error_path"))
}
func loadCustomLogger(log_name string) *zap.Logger {
	return loadLogger(GetConfigString("logger.custom_dir") + log_name + "." + GetConfigString("logger.custom_suffix"))
}
func loadLogger(path string) *zap.Logger {
	if logger_customs == nil {
		logger_customs = map[string]*zap.Logger{}
	}
	if logger_customs[path] != nil {
		return logger_customs[path]
	}
	fmt.Println("init logger " + path)
	logger_customs[path] = newLogger(path)
	return logger_customs[path]
}
func newLogger(path string) *zap.Logger {
	lumberJackLogger := GetLogIoWriter(path)
	writer := zapcore.AddSync(lumberJackLogger)

	encoder_config := zap.NewProductionEncoderConfig()
	encoder_config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder_config.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoder_config)

	logger_level, _ := zapcore.ParseLevel(GetConfigString("logger.level"))
	core := zapcore.NewCore(encoder, writer, logger_level)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func GetLogIoWriter(path string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    GetConfigInt("logger.max_size"),    //100MB
		MaxBackups: GetConfigInt("logger.max_backups"), //保留旧文件的最大个数
		MaxAge:     GetConfigInt("logger.max_age"),     //保留旧文件的最大天数
		Compress:   false,
	}
}
