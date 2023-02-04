package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// It's not recommended to use zap directly cause of dependency injection, if you want to use its function so export it as logger package's function
var Logger *zap.SugaredLogger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./logs/error.log"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006 Jan 02 15:04:05")
	config.EncoderConfig.StacktraceKey = ""

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/error.log",
		MaxSize:    10, //MB
		MaxBackups: 30,
		MaxAge:     30, //days
		Compress:   true,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	logger := zap.New(core)

	Logger = logger.Sugar()

}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
