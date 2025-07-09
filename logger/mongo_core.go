package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MongoCore struct {
	collection *mongo.Collection
	encoder    zapcore.Encoder
	level      zapcore.LevelEnabler
}

func NewMongoCoreFromCollection(coll *mongo.Collection) zapcore.Core {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.LevelKey = "level"
	encoderCfg.NameKey = "logger"
	encoderCfg.CallerKey = "caller"
	encoderCfg.MessageKey = "message"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	return &MongoCore{
		collection: coll,
		encoder:    encoder,
		level:      zapcore.DebugLevel,
	}
}

func (m *MongoCore) Enabled(lvl zapcore.Level) bool {
	return m.level.Enabled(lvl)
}

func (m *MongoCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if m.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, m)
	}
	return checkedEntry
}

func (m *MongoCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *m
	return &clone
}

func (m *MongoCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buf, err := m.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	var logDoc map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logDoc); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = m.collection.InsertOne(ctx, logDoc)
	if err != nil {
		fmt.Println("mongo insertedOne error: ", err)
	}
	return err
}

func (m *MongoCore) Sync() error {
	return nil
}
