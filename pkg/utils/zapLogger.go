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
	// 第一个参数是输出的格式 第二个参数是输出的位置 第三个参数是日志输出级别
	core := zapcore.NewCore(getEncoder(),
		//zapcore.NewMultiWriteSyncer 输出到多个位置 比如 文件 console中
		zapcore.NewMultiWriteSyncer(getWriterSyncer(), zapcore.AddSync(os.Stdout)),
		getLevelPriority())

	logger := zap.New(core).Sugar()
	// 显示行号 如[main.go main:42]
	if conf.Cfg.Zap.ShowLine {
		// 获取 调用的文件, 函数名称, 行号
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

// def 输出日志的格式
func getEncoder() zapcore.Encoder {
	// 参考: zap.NewProductionEncoderConfig()
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

	if conf.Cfg.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// def 日志要输出到什么地方
func getWriterSyncer() zapcore.WriteSyncer {
	// 日志路径
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format("2006-01-02") + ".log"

	// 日志分割
	hook := lumberjack.Logger{
		Filename:   stLogFilePath,           // 日志文件路径，默认 os.TempDir()
		MaxSize:    conf.Cfg.Log.MaxSize,    // 每个日志文件保存500M，默认 100M
		MaxBackups: conf.Cfg.Log.MaxBackups, // 保留3个备份，默认不限
		MaxAge:     conf.Cfg.Log.MaxAge,     // 保留28天，默认不限
		Compress:   conf.Cfg.Log.Compress,   // 是否压缩，默认不压缩
	}
	return zapcore.AddSync(&hook)
}

// 获取日志输出级别
func getLevelPriority() zapcore.LevelEnabler {
	switch conf.Cfg.Zap.Level {
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

// 自定义日志输出时间格式
func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(conf.Cfg.Zap.Prefix + t.Format("2006/01/02 15:04:05"))
}
