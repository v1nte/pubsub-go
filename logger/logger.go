package logger

import (
	"os"

	"github.com/v1nte/pubsub-go/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() error {
	var cores []zapcore.Core

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	cores = append(cores, consoleCore)

	mongoCore := NewMongoCoreFromCollection(database.LogsDB)

	cores = append(cores, mongoCore)

	core := zapcore.NewTee(cores...)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zap.ReplaceGlobals(Log)

	return nil
}
