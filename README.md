# clogger
my first go for logger



可以通过引入，clogger，生成日志；

1. 使用NewFileLogger构造写文件日志方法；参数：文件路径，文件名称，日志级别
2. 使用NewConsoleLogger构造控制台日志方法；参数：日志级别



```go
// 构造文件日志
logger = clogger.NewFileLogger("./", "xxx.log", "debug")
// 构造控制台日志
logger = clogger.NewConsoleLogger("debug")
```

## 使用方法

输出日志，可以采用format格式进行日志输出。

```go
		sb := "管大妈"
		logger.Debug("%s是个好捧哏", sb)
		logger.Info("这是一条Info日志")
		logger.Error("这是一条Error日志")
		logger.Warn("这是一条Warn日志")
		logger.Fatal("这是一条Fatal日志")
```



## 优势

1. 输出日志采用独立的`goroutine`进行输出；
2. 独立的`goroutine`采用绑定的`chan`进行日志任务的接收；
3. 综合前两条，可以将日志输出，异步的提交给独立的goroutine进行，不占用业务程序的处理载荷。

