package main

import (
	"github.com/yz271544/clogger"
)

var logger clogger.Logger

// 一个使用自定义日志库的用户程序
func main() {
	//filePath := `D:\iProject\clogger\logs/`
	//logger = clogger.NewFileLogger(filePath, "xxx.log", "debug", true, true)
	logger = clogger.NewConsoleLogger("debug")

	defer logger.Close()

	/*fmt.Printf("Main GoId %d\n", clogger.Goid())
	sb := "管大妈"
	logger.Debug("1 %s是个好捧哏", sb)
	logger.Info("2 这是一条Info日志")
	logger.Error("3 这是一条Error日志")
	logger.Warn("4 这是一条Warn日志")
	logger.Fatal("5 这是一条Fatal日志")*/

	for {
		sb := "管大妈"
		logger.Debugf("%s是个好捧哏", sb)
		logger.Info("这是一条Info日志")
		logger.Error("这是一条Error日志")
		logger.Warn("这是一条Warn日志")
		logger.Fatal("这是一条Fatal日志")
	}
	//time.Sleep(1 * time.Second)
}
