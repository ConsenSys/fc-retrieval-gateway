package logging

import (
	"fmt"
	"log"
	"sync"
)

const (
	// LogLevelStrNone log level should be used when no logging should be emitted
	LogLevelStrNone = "NONE"
	// LogLevelStrError log level should be used when only Error level logging should be emitted
	LogLevelStrError = "ERROR"
	// LogLevelStrWarn log level should be used when only Error and Warn level logging should be emitted
	LogLevelStrWarn = "WARN"
	// LogLevelStrInfo log level should be used when only Error, Warn, and Info level logging should be emitted
	LogLevelStrInfo = "INFO"
	// LogLevelStrTrace log level should be used when all logging should be emitted
	LogLevelStrTrace = "TRACE"
	
	// LogLevelNone log level should be used when no logging should be emitted
	LogLevelNone = 0
	// LogLevelError log level should be used when only Error level logging should be emitted
	LogLevelError = 1
	// LogLevelWarn log level should be used when only Error and Warn level logging should be emitted
	LogLevelWarn = 2
	// LogLevelInfo log level should be used when only Error, Warn, and Info level logging should be emitted
	LogLevelInfo = 3
	// LogLevelTrace log level should be used when all logging should be emitted
	LogLevelTrace = 4

	logTargetStdOut = "STDOUT"
	logTargetStdOutInt = 0
)


var logLevel = int(LogLevelTrace)
var logTarget = int(logTargetStdOutInt)
var mutex sync.RWMutex


// Init inialises the logging system
func Init() {
	mutex = sync.RWMutex{}
}

// SetLogLevel allows the log level to be specified.
func SetLogLevel(level string) {
	switch (level) {
	case LogLevelStrError:
		logLevel = LogLevelError
	case LogLevelStrWarn:
		logLevel = LogLevelWarn
	case LogLevelStrInfo:
		logLevel = LogLevelInfo
	case LogLevelStrTrace:
		logLevel = LogLevelTrace
	default:
		panic("Unknown Log Level: " + level)
	}
}

// SetLogTarget allows the destination of logs to be specified.
func SetLogTarget(target string) {
	switch (target) {
	case logTargetStdOut:
		logTarget = logTargetStdOutInt
	default:
		panic("Unknown Log Target: " + target)
	}
}

// ErrorEnabled returns true if ERROR log level is enabled.
func ErrorEnabled() bool {
	return LogLevelError <= logLevel
}

// WarnEnabled returns true if WARN log level is enabled.
func WarnEnabled() bool {
	return LogLevelWarn <= logLevel
}

// InfoEnabled returns true if INFO log level is enabled.
func InfoEnabled() bool {
	return LogLevelInfo <= logLevel
}

// TraceEnabled returns true if TRACE log level is enabled.
func TraceEnabled() bool {
	return LogLevelTrace <= logLevel
}

// Error prints out msg to the log target if the log level is ERROR or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Error(msg string, args ...interface{}) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, msg, args...)
	}
}

// Error1 prints out msg to the log target if the log level is ERROR or lower. 
func Error1(err error) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, err.Error())
	}
}

// ErrorAndPanic logs an error and then calls panic with the same message.
func ErrorAndPanic(msg string, args ...interface{}) {
	if (LogLevelError <= logLevel) {
		printf(LogLevelStrError, msg, args...)
	}
	panic(fmt.Sprintf(msg, args...))
}

// Warn prints out msg to the log target if the log level is WARN or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Warn(msg string, args ...interface{}) {
	if (LogLevelWarn <= logLevel) {
		printf(LogLevelStrWarn, msg, args...)
	}
}

// Info prints out msg to the log target if the log level is INFO or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Info(msg string, args ...interface{}) {
	if (LogLevelInfo <= logLevel) {
		printf(LogLevelStrInfo, msg, args...)
	}
}

// Trace prints out msg to the log target if the log level is INFO or lower. 
// msg is interpreted as a format string and args as parameters to the format
// string is there are any args.
func Trace(msg string, args ...interface{}) {
	if (LogLevelTrace <= logLevel) {
		printf(LogLevelStrTrace, msg, args...)
	}
}

func logLevelEnabled(level int) bool {
	return level <= logLevel
}


func printf(level string, msg string, args ...interface{}) {
	var s string
	if (len(args) > 0) {
		s = fmt.Sprintf(msg, args...)
	} else {
		s = msg
	}
	s = level + ": " + s

	switch (logTarget) {
	case logTargetStdOutInt:
		// TODO will need to do better than this in the server so we don't block all of the time!
		mutex.Lock()
		defer mutex.Unlock()
		log.Println(s)
	default:
		panic("Unknown log target")
	}
}