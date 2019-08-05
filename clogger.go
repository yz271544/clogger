package clogger

import "strings"

// 我的日志库文件
// Level 是一个自定义的类型，代表日志级别
type Level uint16

// 定义具体的日志级别
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// 定义一个日志传参体
type Log struct {
	level    Level
	format   string
	args     []interface{}
	fileName string
	line     int
	funcName string
}

// 初始化日志，并获得代码所在文件，行号，函数
func NewLog(level Level, format string, args ...interface{}) *Log {
	fileName, line, funcName := getCallerInfo(3)
	return &Log{
		level:    level,
		format:   format,
		args:     args,
		fileName: fileName,
		line:     line,
		funcName: funcName,
	}
}

// 为Log绑定方法，用于日志处理调度执行
func (l *Log) task(f func(fileName string, line int, funcName string, level Level, format string, args ...interface{})) {
	f(l.funcName, l.line, l.funcName, l.level, l.format, l.args...)
}

// 定义一个Logger接口
type Logger interface {
	Debug(format string, args ...interface{})

	// 方法 info方法
	Info(format string, args ...interface{})

	// 方法 warn方法
	Warn(format string, args ...interface{})

	// 方法 error方法
	Error(format string, args ...interface{})

	// 方法 fatal方法
	Fatal(format string, args ...interface{})

	Close()
}

// 写一个根据传进来的Level，获取对应的字符串
func getLevelStr(level Level) string {
	switch level {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Debug"
	}
}

// 根据用户传入的字符串类型的日志级别，解析出对应的Level
func parseLogLevel(levelStr string) Level {
	levelStr = strings.ToLower(levelStr) // 将字符串转换为全小写
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}
