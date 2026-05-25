package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}

// SetLogLevel changes the global log level.
func SetLogLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// Info logs informational messages at the Info level.
func Info(args ...interface{})                 { logger.Info(args...) }
// Infof logs formatted informational messages at the Info level.
func Infof(format string, args ...interface{}) { logger.Infof(format, args...) }
// Infoln logs informational messages with a newline at the Info level.
func Infoln(args ...interface{})               { logger.Infoln(args...) }
// Warn logs warning messages at the Warn level.
func Warn(args ...interface{})                 { logger.Warn(args...) }
// Warnf logs formatted warning messages at the Warn level.
func Warnf(format string, args ...interface{}) { logger.Warnf(format, args...) }
// Warnln logs warning messages with a newline at the Warn level.
func Warnln(args ...interface{})               { logger.Warnln(args...) }
// Error logs error messages at the Error level.
func Error(args ...interface{})                { logger.Error(args...) }
// Errorf logs formatted error messages at the Error level.
func Errorf(format string, args ...interface{}){ logger.Errorf(format, args...) }
// Errorln logs error messages with a newline at the Error level.
func Errorln(args ...interface{})              { logger.Errorln(args...) }
// Fatal logs messages at the Fatal level and exits.
func Fatal(args ...interface{})                { logger.Fatal(args...) }
// Fatalf logs formatted messages at the Fatal level and exits.
func Fatalf(format string, args ...interface{}){ logger.Fatalf(format, args...) }
// Fatalln logs messages with a newline at the Fatal level and exits.
func Fatalln(args ...interface{})              { logger.Fatalln(args...) }
