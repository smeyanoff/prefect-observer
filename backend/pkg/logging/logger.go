package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init() {
	// Настраиваем encoder для вывода логов
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	// Настраиваем вывод логов в стандартный поток (os.Stdout)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	// Настраиваем вывод логов в файл
	file, err := os.OpenFile("../backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("can not open log file")
	}
	fileWriteSyncer := zapcore.AddSync(file)

	// Устанавливаем уровень логирования
	logLevel := zapcore.DebugLevel

	// Создаем core с множественными выводами (в консоль и файл)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriteSyncer, logLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriteSyncer, logLevel),
	)

	// Создаем логгер с core и добавляем вызов Stacktrace для Error уровня и выше
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func WrapError(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", prefix, err)
}

// Debug выводит отладочные сообщения
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info выводит информационные сообщения
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn выводит предупреждения
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error выводит сообщения об ошибках
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal выводит сообщения об ошибках с завершением программы
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync гарантирует, что все буферы логов будут сброшены
func Sync() {
	_ = Logger.Sync()
}
