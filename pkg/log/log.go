package log

import (
	"douyin/conf"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOGGER_KEY = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLog(cfg *conf.Config) *Logger {
	return initZap(cfg)
}

func initZap(cfg *conf.Config) *Logger {
	// 第一个参数是输出的格式 第二个参数是输出的位置 第三个参数是日志输出级别
	core := zapcore.NewCore(getEncoder(cfg),
		//zapcore.NewMultiWriteSyncer 输出到多个位置 比如 文件 console中
		zapcore.NewMultiWriteSyncer(getWriterSyncer(cfg), zapcore.AddSync(os.Stdout)),
		getLevelPriority(cfg))

	// if conf.GetString("env") != "prod" {
	// 	return &Logger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	// }

	if cfg.Log.ShowLine {
		return &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
	return &Logger{zap.New(core, zap.AddStacktrace(zap.ErrorLevel))}
}

// NewContext Adds a field to the specified context
func (l *Logger) NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(LOGGER_KEY, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(ctx *gin.Context) *Logger {
	if ctx == nil {
		return l
	}
	zl, _ := ctx.Get(LOGGER_KEY)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}

// def 输出日志的格式
func getEncoder(cfg *conf.Config) zapcore.Encoder {
	// 参考: zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(cfg.Log.Prefix + t.Format("2006/01/02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	if cfg.Log.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// def 日志要输出到什么地方
func getWriterSyncer(cfg *conf.Config) zapcore.WriteSyncer {
	// 日志路径
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + cfg.Log.Directory + stSeparator + time.Now().Format("2006-01-02") + ".log"

	// 日志分割
	hook := lumberjack.Logger{
		Filename:   stLogFilePath,      // 日志文件路径，默认 os.TempDir()
		MaxSize:    cfg.Log.MaxSize,    // 每个日志文件保存500M，默认 100M
		MaxBackups: cfg.Log.MaxBackups, // 保留3个备份，默认不限
		MaxAge:     cfg.Log.MaxAge,     // 保留28天，默认不限
		Compress:   cfg.Log.Compress,   // 是否压缩，默认不压缩
	}
	return zapcore.AddSync(&hook)
}

// 获取日志输出级别
func getLevelPriority(cfg *conf.Config) zapcore.LevelEnabler {
	switch cfg.Log.Level {
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
