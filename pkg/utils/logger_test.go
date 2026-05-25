package utils

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// saveLogger saves and restores the package-level logger around a test.
func saveLogger(t *testing.T) {
	t.Helper()
	orig := logger
	t.Cleanup(func() { logger = orig })
}

// captureOutput replaces the logger output with a buffer, runs fn, and returns the output.
// It restores the original output and level after the test.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	origOut := logger.Out
	origLevel := logger.Level
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel)
	defer func() {
		logger.SetOutput(origOut)
		logger.SetLevel(origLevel)
	}()
	fn()
	return buf.String()
}

func TestLoggerInitialized(t *testing.T) {
	assert.NotNil(t, logger, "Logger should be initialized by init()")
}

func TestLoggerFormat(t *testing.T) {
	assert.NotNil(t, logger)

	formatter := logger.Formatter
	textFormatter, ok := formatter.(*logrus.TextFormatter)
	assert.True(t, ok, "Logger formatter should be TextFormatter")
	assert.True(t, textFormatter.FullTimestamp, "FullTimestamp should be enabled")
	assert.Equal(t, "2006-01-02 15:04:05", textFormatter.TimestampFormat,
		"TimestampFormat should match expected format")
}

func TestLoggerLevel(t *testing.T) {
	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.Level, "Logger level should be InfoLevel")
}

func TestLoggerOutput(t *testing.T) {
	assert.NotNil(t, logger)
	assert.Equal(t, os.Stdout, logger.Out, "Logger output should be stdout")
}

func TestSetLogLevel(t *testing.T) {
	origLevel := logger.Level
	defer func() { logger.SetLevel(origLevel) }()

	SetLogLevel(logrus.DebugLevel)
	assert.Equal(t, logrus.DebugLevel, logger.Level)

	SetLogLevel(logrus.WarnLevel)
	assert.Equal(t, logrus.WarnLevel, logger.Level)
}

func TestInfo(t *testing.T) {
	output := captureOutput(t, func() {
		Info("test info message")
	})
	assert.Contains(t, output, "test info message")
	assert.Contains(t, output, "level=info")
}

func TestInfof(t *testing.T) {
	output := captureOutput(t, func() {
		Infof("formatted %s %d", "info", 42)
	})
	assert.Contains(t, output, "formatted info 42")
	assert.Contains(t, output, "level=info")
}

func TestInfoln(t *testing.T) {
	output := captureOutput(t, func() {
		Infoln("infoln message")
	})
	assert.Contains(t, output, "infoln message")
}

func TestWarn(t *testing.T) {
	output := captureOutput(t, func() {
		Warn("test warn message")
	})
	assert.Contains(t, output, "test warn message")
	assert.Contains(t, output, "level=warning")
}

func TestWarnf(t *testing.T) {
	output := captureOutput(t, func() {
		Warnf("formatted %s %d", "warn", 99)
	})
	assert.Contains(t, output, "formatted warn 99")
	assert.Contains(t, output, "level=warning")
}

func TestWarnln(t *testing.T) {
	output := captureOutput(t, func() {
		Warnln("warnln message")
	})
	assert.Contains(t, output, "warnln message")
}

func TestError(t *testing.T) {
	output := captureOutput(t, func() {
		Error("test error message")
	})
	assert.Contains(t, output, "test error message")
	assert.Contains(t, output, "level=error")
}

func TestErrorf(t *testing.T) {
	output := captureOutput(t, func() {
		Errorf("formatted %s %d", "error", 7)
	})
	assert.Contains(t, output, "formatted error 7")
	assert.Contains(t, output, "level=error")
}

func TestErrorln(t *testing.T) {
	output := captureOutput(t, func() {
		Errorln("errorln message")
	})
	assert.Contains(t, output, "errorln message")
	assert.Contains(t, output, "level=error")
}

// testExitHook captures the message passed to logrus.ExitFunc.
// This lets us test Fatal/Fatalf/Fatalln without actually calling os.Exit.
func setupExitHook(t *testing.T) (chan int, func()) {
	t.Helper()
	saveLogger(t)

	ch := make(chan int, 1)
	origExit := logger.ExitFunc
	logger.ExitFunc = func(code int) { ch <- code }

	return ch, func() { logger.ExitFunc = origExit }
}

func TestFatal(t *testing.T) {
	exitCh, restore := setupExitHook(t)
	defer restore()

	output := captureOutput(t, func() {
		Fatal("fatal message")
	})
	assert.Contains(t, output, "fatal message")
	assert.Contains(t, output, "level=fatal")
	select {
	case code := <-exitCh:
		assert.Equal(t, 1, code)
	default:
		t.Fatal("expected ExitFunc to be called")
	}
}

func TestFatalf(t *testing.T) {
	exitCh, restore := setupExitHook(t)
	defer restore()

	output := captureOutput(t, func() {
		Fatalf("formatted %s %d", "fatal", 1)
	})
	assert.Contains(t, output, "formatted fatal 1")
	assert.Contains(t, output, "level=fatal")
	select {
	case code := <-exitCh:
		assert.Equal(t, 1, code)
	default:
		t.Fatal("expected ExitFunc to be called")
	}
}

func TestFatalln(t *testing.T) {
	exitCh, restore := setupExitHook(t)
	defer restore()

	output := captureOutput(t, func() {
		Fatalln("fatalln message")
	})
	assert.Contains(t, output, "fatalln message")
	assert.Contains(t, output, "level=fatal")
	select {
	case code := <-exitCh:
		assert.Equal(t, 1, code)
	default:
		t.Fatal("expected ExitFunc to be called")
	}
}
