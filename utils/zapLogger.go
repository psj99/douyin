package utils

import (
	"douyin/conf"

	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ZapLogger *zap.SugaredLogger

func InitLogger() {
	ZapLogger = newZapLogger()
}

func newZapLogger() *zap.SugaredLogger {
	core := zapcore.NewCore(getEncoder(), // 输出格式
		zapcore.NewMultiWriteSyncer(getWriterSyncer(), zapcore.AddSync(os.Stdout)), // 输出位置+同时输出到控制台
		getLevelPriority()) // 输出级别

	logger := zap.New(core).Sugar()

	if conf.Cfg.Log.ShowLine {
		logger = logger.WithOptions(zap.AddCaller()) // 显示[调用的文件 函数名称:行号] 如[main.go main:42]
	}

	return logger
}

// 输出格式
func getEncoder() zapcore.Encoder {
	// 参考 zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 输出时间格式
func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(conf.Cfg.Log.Prefix + t.Format("2006/01/02 15:04:05"))
}

// 输出位置
func getWriterSyncer() zapcore.WriteSyncer {
	// 日志路径
	stLogFilePath := filepath.Join(conf.Cfg.Log.Path, (time.Now().Format("2006-01-02") + ".log")) // 拼接路径
	stLogFilePath = filepath.Clean(stLogFilePath)                                                 // 清理路径语法
	stLogFilePath, err := filepath.Abs(stLogFilePath)                                             // 转换为绝对路径
	if err != nil {
		panic(err)
	}

	// 日志分割
	hook := lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    conf.Cfg.Log.MaxSize,
		MaxBackups: conf.Cfg.Log.MaxBackups,
		MaxAge:     conf.Cfg.Log.MaxAge,
		Compress:   conf.Cfg.Log.Compress,
	}
	return zapcore.AddSync(&hook)
}

// 输出级别
func getLevelPriority() zapcore.LevelEnabler {
	switch conf.Cfg.Log.Level {
	case "debug", "Debug":
		return zap.DebugLevel
	case "info", "Info":
		return zap.InfoLevel
	case "warn", "Warn":
		return zap.WarnLevel
	case "error", "Error":
		return zap.ErrorLevel
	case "dpanic", "DPanic":
		return zap.DPanicLevel
	case "panic", "Panic":
		return zap.PanicLevel
	case "fatal", "Fatal":
		return zap.FatalLevel
	}
	return zap.InfoLevel
}
