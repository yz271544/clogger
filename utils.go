package clogger

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
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
