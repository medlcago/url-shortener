package logger

import (
	"github.com/sirupsen/logrus"
	"url-shortener/config"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
}

type LogrusLogger struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewLogrusLogger(cfg *config.Config) *LogrusLogger {
	l := &LogrusLogger{cfg: cfg}
	l.initLogger()
	return l
}

func (l *LogrusLogger) getLoggerLevel() logrus.Level {
	level, err := logrus.ParseLevel(l.cfg.Logger.Level)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func (l *LogrusLogger) initLogger() {
	logrusLog := logrus.New()
	if l.cfg.Logger.Encoding == "json" {
		logrusLog.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrusLog.SetFormatter(&logrus.TextFormatter{})
	}
	logLever := l.getLoggerLevel()
	logrusLog.SetLevel(logLever)

	l.logger = logrusLog
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *LogrusLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *LogrusLogger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *LogrusLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}
