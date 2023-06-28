package logger

import (
	"cess-faucet/config"
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	InfoLogger *zap.Logger
	ErrLogger  *zap.Logger
)

func LoggerInit() {
	_, err := os.Stat(config.LogfilePathPrefix)
	if err != nil {
		err = os.MkdirAll(config.LogfilePathPrefix, os.ModePerm)
		if err != nil {
			config.LogfilePathPrefix = "./log/"
		}
	}
	initInfoLogger()
	initErrLogger()
}

// info log
func initInfoLogger() {
	infologpath := config.LogfilePathPrefix + "info.log"
	hook := lumberjack.Logger{
		Filename:   infologpath,
		MaxSize:    10,
		MaxAge:     360,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		TimeKey:      "time",
		CallerKey:    "file",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   formatEncodeTime,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	InfoLogger = zap.New(core, caller, development)
	InfoLogger.Sugar().Infof("The service has started and created a log file in the %v", infologpath)
}

// error log
func initErrLogger() {
	errlogpath := config.LogfilePathPrefix + "error.log"
	hook := lumberjack.Logger{
		Filename:   errlogpath,
		MaxSize:    10,
		MaxAge:     360,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		TimeKey:      "time",
		CallerKey:    "file",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   formatEncodeTime,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.ErrorLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	ErrLogger = zap.New(core, caller, development)
	ErrLogger.Sugar().Errorf("The service has started and created a log file in the %v", errlogpath)
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}
