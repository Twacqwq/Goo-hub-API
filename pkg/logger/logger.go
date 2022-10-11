// Package logger 处理日志相关逻辑
package logger

import (
	"fmt"
	"os"
	"strings"
	"thub/pkg/app"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局Logger对象
var Logger *zap.Logger

// 日志初始化
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType, level string) {
	// 获取写入介质
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)

	//设置日志等级
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误, 日志级别设置有误。")
	}

	// 初始化core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 初始化Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error时才会显示 stacktrace
	)

	// 将自定义的Logger替换为全局的Logger
	// zap.L().Fatal() 调用时, 就会使用我们自定义的Logger
	zap.ReplaceGlobals(Logger)
}

// 设置日志存储格式
func getEncoder() zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志结尾添加"\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写 如: ERROR INFO
		EncodeTime:     customTimeEncoder,              // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间 以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式

	}

	// 本地环境配置
	if app.IsLocal() {
		// 终端输出关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 线上环境使用JSON编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 日志记录介质 Thub使用两种介质 os.Stdout & file  返回一个写入器(介质)
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 如果配置了按照日期记录日志文件
	if logType == "daily" {
		logname := time.Now().Format("2006-01-02.log") //Format() 参数代表模板格式
		filename = strings.ReplaceAll(filename, "logs.log", logname)
	}

	// 滚动日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	// 配置输出介质
	if app.IsLocal() {
		// 本地开发终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产环境只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}
}
