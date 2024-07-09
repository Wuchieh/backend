package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"string_backend_0001/internal/conf"
	"strings"
	"time"
)

var (
	logFile    *os.File
	logger     *log.Logger
	level      int
	logPath    string
	consoleLog bool
	lastDate   string
)

func GetLogger() *log.Logger {
	return logger
}

func Init() error {
	logConf := conf.Conf.Log
	logPath = logConf.Path

	switch strings.ToLower(logConf.Level) {
	case "debug":
		level = LDebug
		consoleLog = true
	case "info":
		level = LInfo
	case "warn":
		level = LWarn
	case "error":
		level = LError
	default:
		return errors.New("log level error")
	}

	if !strings.HasSuffix(logPath, ".log") {
		return errors.New("log file error, path must end with .log")
	}

	err := openLogFile()
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	var writers []io.Writer
	writers = append(writers, logFile)
	if consoleLog {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)

	logger = log.New(multiWriter, "", 0)
	lastDate = time.Now().Format(time.DateOnly)

	return nil
}

func openLogFile() error {
	var err error
	err = os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
	if err != nil {
		return err
	}
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return err
}

func checkRotate() {
	if logFile == nil {
		return
	}

	currentDate := time.Now().Format(time.DateOnly)
	if currentDate != lastDate {
		// Close the current log file
		err := logFile.Close()
		if err != nil {
			Error("failed to close log file")
			return
		}

		// Rename the old log file
		oldPath := logPath
		newPath := strings.TrimSuffix(oldPath, ".log") + "-" + lastDate + ".log"
		err = os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("Failed to rename log file: %v\n", err)
		}

		// Open a new log file
		err = openLogFile()
		if err != nil {
			fmt.Printf("Failed to open new log file: %v\n", err)
		}

		// Update the lastDate
		lastDate = currentDate

		// Recreate the logger with the new file
		var writers []io.Writer
		writers = append(writers, logFile)
		if consoleLog {
			writers = append(writers, os.Stdout)
		}
		multiWriter := io.MultiWriter(writers...)
		logger = log.New(multiWriter, "", 0)
	}
}

func Info(format string, v ...interface{}) {
	if level <= LInfo {
		output(LInfo, format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if level <= LDebug {
		output(LDebug, format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if level <= LWarn {
		output(LWarn, format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if level <= LError {
		output(LError, format, v...)
	}
}

func output(logLevel int, format string, v ...interface{}) {
	checkRotate()

	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(2)
	msg := fmt.Sprintf(format, v...)

	var levelStr string
	switch logLevel {
	case LDebug:
		levelStr = "DEBUG"
	case LInfo:
		levelStr = "INFO"
	case LWarn:
		levelStr = "WARN"
	case LError:
		levelStr = "ERROR"
	default:
		levelStr = "OTHER"
	}

	logMsg := fmt.Sprintf("[%s]\t| %s | %s:%d | %s", levelStr, now, filepath.Base(file), line, msg)
	logger.Println(logMsg)
}

func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}
