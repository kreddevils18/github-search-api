package log

import (
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ILogger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logger struct {
	log *zap.SugaredLogger
}

func NewLogger(conf *config.Config) (*Logger, error) {
	var encoderConfig zapcore.EncoderConfig
	if os.Getenv("NODE_ENV") != "production" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "CALLER"
	encoderConfig.TimeKey = "TIME"
	encoderConfig.LevelKey = "LEVEL"
	encoderConfig.MessageKey = "MESSAGE"
	encoderConfig.StacktraceKey = ""

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	logFile, _ := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, os.Stdout, defaultLogLevel))

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()

	return &Logger{log: sugar}, nil
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.log.Infof(template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.log.Fatalf(template, args...)
}
