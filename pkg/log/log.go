package log

import (
	"context"
	"errors"
	"os"
	"tinyurl/pkg/config"
	"tinyurl/pkg/trace"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
}

const (
	StdOutput  = "std"
	FileOutput = "file"
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

var colors = map[color.Attribute]*color.Color{
	color.FgGreen:   color.New(color.FgGreen),
	color.FgHiWhite: color.New(color.FgHiWhite),
	color.FgYellow:  color.New(color.FgYellow),
	color.FgRed:     color.New(color.FgRed, color.Underline),
	color.FgHiRed:   color.New(color.FgHiRed, color.Underline, color.Bold),
	color.FgBlue:    color.New(color.FgBlue),
}

func New(c config.Log) (*Logger, error) {
	level := getLevel(c.Level)
	var writer zapcore.WriteSyncer
	if c.Output == StdOutput {
		writer = getStdWriter()
	} else if c.Output == FileOutput {
		writer = getFileWriter(c.Rotate)
	} else {
		return nil, errors.New("output must in [std, file]")
	}
	encoder := getJSONEncoder()
	var opts []zap.Option
	if c.Development {
		encoder = getConsoleEncoder()
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	}
	core := zapcore.NewCore(encoder, writer, level)
	return &Logger{zap.New(core, opts...)}, nil
}

func getLevel(level string) zapcore.Level {
	if l, ok := levelMap[level]; ok {
		return l
	}
	return zapcore.InfoLevel
}

func getJSONEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	conf.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewJSONEncoder(conf)
}

func getConsoleEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	conf.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewConsoleEncoder(conf)
}

func getFileWriter(c config.Rotate) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		LocalTime:  c.LocalTime,
		Compress:   c.Compress,
	}
	return zapcore.AddSync(logger)
}

func getStdWriter() zapcore.WriteSyncer {
	return os.Stdout
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgBlue].Sprintf("%s", msg)
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgGreen].Sprintf("%s", msg)
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgYellow].Sprintf("%s", msg)
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgRed].Sprintf("%s", msg)
	l.logger.Error(msg, fields...)
}

func (l *Logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgHiRed].Sprintf("%s", msg)
	l.logger.DPanic(msg, fields...)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgHiRed].Sprintf("%s", msg)
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getTrace(ctx))
	msg = colors[color.FgHiRed].Sprintf("%s", msg)
	l.logger.Fatal(msg, fields...)
}

func getTrace(ctx context.Context) zapcore.Field {
	if ctx == nil {
		return zap.Skip()
	}
	return zap.String("trace", trace.Trace(ctx))
}
