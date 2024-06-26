package log

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_logger *zap.Logger
	aLevel  zap.AtomicLevel
)

const (
	// FormatText format log text
	FormatText = "text"
	// FormatJSON format log json
	FormatJSON = "json"
)

// type Level uint

// 日志配置
type logConfig struct {
	LogPath    string
	LogLevel   string
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Format     string
}

func init() {
	InitLog()
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logPath string, maxSize int, compress bool) io.Writer {
	if logPath == "" || logPath == "-" {
		log.Println("output logs to stdout")
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename: logPath,
		MaxSize:  maxSize,
		Compress: compress,
	}
}

func newZapEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}
func newLoggerCore(log *logConfig) zapcore.Core {
	hook := newLogWriter(log.LogPath, log.MaxSize, log.Compress)

	encoderConfig := newZapEncoder()

	aLevel = zap.NewAtomicLevelAt(getZapLevel(log.LogLevel))

	var encoder zapcore.Encoder
	if log.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(hook),
		aLevel,
	)
	return core
}

func newLoggerOptions() []zap.Option {
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	development := zap.Development()
	options := []zap.Option{
		caller,
		callerSkip,
		development,
	}
	return options
}

// Option function option
type Option func(*logConfig)

// Path set logPath
// if is zero will print,or write file
func Path(logPath string) Option {
	return func(logCfg *logConfig) {
		logCfg.LogPath = logPath
	}
}

// Compress compress log
func Compress(compress bool) Option {
	return func(logCfg *logConfig) {
		logCfg.Compress = compress
	}
}

// Level set log level default info
func Level(level string) Option {
	return func(logCfg *logConfig) {
		logCfg.LogLevel = level
	}
}

// MaxSize Log Max Size
func MaxSize(size int) Option {
	return func(logCfg *logConfig) {
		logCfg.MaxSize = size
	}
}

// MaxAge log store day
func MaxAge(age int) Option {
	return func(logCfg *logConfig) {
		logCfg.MaxAge = age
	}
}

// MaxBackups total store log
func MaxBackups(backup int) Option {
	return func(logCfg *logConfig) {
		logCfg.MaxBackups = backup
	}
}

// Format log json or text
func Format(format string) Option {
	return func(logCfg *logConfig) {
		if format == FormatJSON {
			logCfg.Format = FormatJSON
		} else {
			logCfg.Format = FormatText
		}

	}
}

func defaultOption() *logConfig {
	return &logConfig{
		LogPath:    "",
		MaxSize:    20,
		Compress:   true,
		MaxAge:     7,
		MaxBackups: 7,
		LogLevel:   "debug",
		Format:     FormatText,
	}
}

// InitLog conf
func InitLog(opts ...Option) {
	if _logger != nil && len(opts) == 0 {
		return
	}
	logCfg := defaultOption()
	for _, opt := range opts {
		opt(logCfg)
	}
	core := newLoggerCore(logCfg)

	zapOpts := newLoggerOptions()
	_logger = zap.New(core, zapOpts...)
}

func SetLevel(level string) {
	aLevel.SetLevel(getZapLevel(level))
}

// Debug output log
func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

// Info output log
func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

// Warn output log
func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

// Error output log
func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

// Panic output panic
func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

// Fatal output log
func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
