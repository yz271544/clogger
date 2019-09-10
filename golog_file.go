package clogger

import (
	"fmt"
	"github.com/kataras/golog"
	"os"
	"strings"
	"sync"
	"time"
)

// 往文件里写日志

type GologFileLogger struct {
	logger   *golog.Logger
	level    Level
	fileName string
	filePath string

	LogFilePrefix string
	CurrentTime   func() string

	IsRotateByDay  bool
	IsRotateBySize bool
	RotateBySize   int64
	RotateNum      int

	file *os.File

	jobChan chan *Log
	sync.Once
}

// 构造函数
func NewGologFileLogger(filePath, logFilePrefix, level string, rotateByDay, rotateBySize bool) *GologFileLogger {
	logLevel := parseLogLevel(level)
	jobChan := make(chan *Log, 100)
	fl := GologFileLogger{
		logger:   golog.New(),
		level:    logLevel,
		fileName: generateLogFileFullName(filePath, logFilePrefix, 0),
		filePath: filePath,

		LogFilePrefix: logFilePrefix,
		CurrentTime:   func() string { return time.Now().Format(LOG_TIME_FORMAT_FILE_SEGMENT) },

		IsRotateByDay:  rotateByDay,
		IsRotateBySize: rotateBySize,
		RotateBySize:   10 * 1024 * 1024,
		RotateNum:      0,

		jobChan: jobChan,
	}
	fl.initFile()  // 根据上面的文件路径和文件名打开日志文件，把文件句柄赋值给结构体
	go fl.doTask() // 在构造函数中，创建goroutine，用于日志处理调度
	return &fl
}

// 从chan中获取Log实体，并对调用日志处理函数
func (c *GologFileLogger) doTask() {
	for log := range c.jobChan {
		log.task(c.processLog)
	}
}

// 将指定的日志文件打开，赋值给结构体
func (f *GologFileLogger) initFile() {
	//logName := path.Join(f.filePath, f.fileName)

	// 打开日志文件
	fileObj, err := os.OpenFile(f.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("打开日志文件%s失败, %v", f.fileName, err))
	}
	f.file = fileObj

	// 打开错误文件
	//reg := regexp.MustCompile(`(.*)\.log$`)
	//prefix := reg.ReplaceAllString(f.fileName, "$1")
	//errLogFileName := fmt.Sprintf("%s.error.log", prefix)
	//errLogFullFileName := path.Join(f.filePath, errLogFileName)
	//errFileObj, err := os.OpenFile(errLogFullFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	//if err != nil {
	//    panic(fmt.Errorf("打开错误日志文件%s失败, %v", errLogFullFileName, err))
	//}
	//f.errFile = errFileObj

}

// 检查是否要拆分
func (f *GologFileLogger) checkSplit(file *os.File) bool {
	info, _ := file.Stat()
	/*info, err := file.Stat()
	  if err != nil {
	  	return false
	  }*/
	fileSize := info.Size()
	return fileSize >= f.RotateBySize // 当传进来的日志文件大小超过maxSize，就返回true
}

// 封装一个切分日志文件的方法
func (f *GologFileLogger) splitLogFile(fileObj *os.File) *os.File {
	// 切分文件
	toCloseFileName := fileObj.Name() // 这里可以获取到文件的全路径
	toBackupFileName := fmt.Sprintf("%s_%v.bak", toCloseFileName, time.Now().Unix())
	// 1. 把原来的文件关闭
	fileObj.Close()
	// 2. 备份原来的文件
	os.Rename(toCloseFileName, toBackupFileName)
	// 3. 新建一个文件
	newFileObj, err := os.OpenFile(toCloseFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("日志文件创建失败:%s, %v", toCloseFileName, err))
	}
	//f.file = newFileObj // 错误日志应该是写 f.errFile = newFileObj
	return newFileObj
}

// 将公用的记录日志的功能封装成一个单独的方法
func (f *GologFileLogger) processLog(fileName string, line int, funcName string, level Level, format string, args ...interface{}) {
	if f.level > level {
		return
	}

	gologger := f.logger

	//gologger.AddOutput(f.file)

	gologger.SetOutput(f.file)

	//f.file.Write()
	//msg := fmt.Sprintf(format, args...) // 得到用户要记录的日志

	gologLevel := golog.ParseLevel(strings.ToLower(getLevelStr(level)))

	gologger.Logf(gologLevel, format, args...)
	//f.logger.pr

	// 日志格式: [时间][文件:行号][函数名][日志级别]日志信息
	//nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	//fileName, line, funcName := getCallerInfo(3)
	//logLevelStr := getLevelStr(level)
	//logMsg := fmt.Sprintf("[%d][%s][%s:%d][%s][%s]%s", Goid(), nowStr, fileName, line, funcName, logLevelStr, msg)

	// 往文件里写之前，要进行检查
	// 检查当前日志文件的大小是否超过了maxSize
	if f.checkSplit(f.file) {
		f.file = f.splitLogFile(f.file)
	}

	//fmt.Fprintln(f.file, logMsg) // 利用fmt包将msg写入f.file文件中

	// 如果是error或者fatal级别的日志，还要记录到f.errFile
	//if level >= ErrorLevel {
	//    if f.checkSizeSplit(f.errFile) {
	//        f.errFile = f.splitLogFile(f.errFile)
	//    }
	//    fmt.Fprintln(f.errFile, logMsg)
	//}
}

// 方法 debug方法
func (f *GologFileLogger) Debug(format string, args ...interface{}) {
	//f.log(DebugLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: DebugLevel, format: format, args: args}
	log := NewLog(DebugLevel, format, args...)
	f.jobChan <- log
}

// 方法 info方法
func (f *GologFileLogger) Info(format string, args ...interface{}) {
	//f.log(InfoLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: InfoLevel, format: format, args: args}
	log := NewLog(InfoLevel, format, args...)
	f.jobChan <- log
}

// 方法 warn方法
func (f *GologFileLogger) Warn(format string, args ...interface{}) {
	//f.log(WarnLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: WarnLevel, format: format, args: args}
	log := NewLog(WarnLevel, format, args...)
	f.jobChan <- log
}

// 方法 error方法
func (f *GologFileLogger) Error(format string, args ...interface{}) {
	//f.log(ErrorLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: ErrorLevel, format: format, args: args}
	log := NewLog(ErrorLevel, format, args...)
	f.jobChan <- log
}

// 方法 fatal方法
func (f *GologFileLogger) Fatal(format string, args ...interface{}) {
	//f.log(FatalLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: FatalLevel, format: format, args: args}
	log := NewLog(FatalLevel, format, args...)
	f.jobChan <- log
}

// Close()
func (f *GologFileLogger) Close() {
	close(f.jobChan)
	f.file.Close()
	//f.errFile.Close()
}
