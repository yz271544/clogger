package clogger

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 存放一些共用的工具函数

func getCallerInfo(skip int) (fileName string, line int, funcName string) {
	pc, fileFullName, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	// 从fileFullName中剥离出文件名
	fileName = path.Base(fileFullName)
	// 根据pc拿到函数名
	funcName = runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)
	return
}

/**
通过日志路径、日志文件前缀、日志文件后缀、循环数字（必须大于0）
*/
func generateLogFileFullName(logDir, logFilePrefix string, rotateNum int) string {
	var (
		rotateNumStr = ""
	)
	today := time.Now().Format(LOG_TIME_FORMAT_FILE_SEGMENT)

	if rotateNum > 0 {
		rotateNumStr = fmt.Sprintf("-%d", rotateNum)
	}

	fileName := fmt.Sprintf("%s-%s%s%s", logFilePrefix, today, rotateNumStr, ".log")
	return path.Join(logDir, fileName)
}

/**
通过日志路径、日志文件前缀、日志文件后缀、循环数字（必须大于0）
*/
func generateLogFileName(logFilePrefix string) string {
	today := time.Now().Format(LOG_TIME_FORMAT_FILE_SEGMENT)
	return fmt.Sprintf("%s-%s%s", logFilePrefix, today, ".log")
}

func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
