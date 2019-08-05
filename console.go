package clogger

import (
	"fmt"
	"os"
	"time"
)

// 往终端打印日志

// ConsoleLogger 是一个终端日志
type ConsoleLogger struct {
	level   Level
	file    *os.File
	jobChan chan *Log
}

// 构造函数
func NewConsoleLogger(level string) *ConsoleLogger {
	logLevel := parseLogLevel(level)
	jobChan := make(chan *Log, 100)
	cl := ConsoleLogger{
		level:   logLevel,
		file:    os.Stdout,
		jobChan: jobChan,
	}
	// 在构造函数中，创建goroutine，用于日志处理调度
	go cl.doTask()

	return &cl
}

func (c *ConsoleLogger) doTask() {
	for log := range c.jobChan {
		log.task(c.processLog)
	}
}

// 将公用的记录日志的功能封装成一个单独的方法
func (c *ConsoleLogger) processLog(fileName string, line int, funcName string, level Level, format string, args ...interface{}) {
	if c.level > level {
		return
	}
	//f.file.Write()
	msg := fmt.Sprintf(format, args...) // 得到用户要记录的日志
	// 日志格式: [时间][文件:行号][函数名][日志级别]日志信息
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	//fileName, line, funcName := getCallerInfo(2)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%d][%s][%s:%d][%s][%s]%s", Goid(), nowStr, fileName, line, funcName, logLevelStr, msg)
	fmt.Fprintln(c.file, logMsg) // 利用fmt包将msg写入f.file文件中
}

// 方法 debug方法
func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	//c.log(DebugLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: DebugLevel, format: format, args: args}
	log := NewLog(DebugLevel, format, args...)
	c.jobChan <- log
}

// 方法 info方法
func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	//c.log(InfoLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: InfoLevel, format: format, args: args}
	log := NewLog(InfoLevel, format, args...)
	c.jobChan <- log
}

// 方法 warn方法
func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	//c.log(WarnLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: WarnLevel, format: format, args: args}
	log := NewLog(WarnLevel, format, args...)
	c.jobChan <- log
}

// 方法 error方法
func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	//c.log(ErrorLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: ErrorLevel, format: format, args: args}
	log := NewLog(ErrorLevel, format, args...)
	c.jobChan <- log
}

// 方法 fatal方法
func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	//c.log(FatalLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: FatalLevel, format: format, args: args}
	log := NewLog(FatalLevel, format, args...)
	c.jobChan <- log
}

// Close 终端标准输出不需要关闭
func (c *ConsoleLogger) Close() {
	close(c.jobChan)
	//c.file.Close()
}
