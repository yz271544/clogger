package clogger

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	FILE_FLUSH_SIZE = 4 * 1024
)

type fileType uint16

const (
	LOG_FILE fileType = iota
	ERR_FILE
)

// 往文件里写日志

type FileLogger struct {
	level      Level
	filePrefix string
	filePath   string

	fileName string

	file    *os.File
	errFile *os.File

	IsRotateByTime bool
	IsRotateBySize bool

	RotateBySize        int64
	RotateLogNum        int
	RotateErrNum        int
	IsMakeRotateNumZero bool

	jobChan     chan *Log
	CurrentTime func() string
	sync.Mutex
}

// 构造函数
func NewFileLogger(filePath, logFilePrefix, level string, rotateByTime, rotateBySize bool) *FileLogger {
	logLevel := parseLogLevel(level)
	jobChan := make(chan *Log, 100)
	fl := FileLogger{
		level:      logLevel,
		filePrefix: logFilePrefix,
		filePath:   filePath,

		fileName: generateLogFileName(logFilePrefix),

		IsRotateByTime:      rotateByTime,
		IsRotateBySize:      rotateBySize,
		RotateBySize:        10 * 1024 * 1024,
		RotateLogNum:        1,
		RotateErrNum:        1,
		IsMakeRotateNumZero: false,

		jobChan: jobChan,
		CurrentTime: func() string {
			return time.Now().Format(LOG_TIME_FORMAT)
		},
	}
	fl.initFile()  // 根据上面的文件路径和文件名打开日志文件，把文件句柄赋值给结构体
	go fl.doTask() // 在构造函数中，创建goroutine，用于日志处理调度
	//go fl.watchFile()
	return &fl
}

// 从chan中获取Log实体，并对调用日志处理函数
func (c *FileLogger) doTask() {
	/*for log := range c.jobChan {
		log.task(c.processLog)
	}*/

	for {
		select {
		case log := <-c.jobChan:
			log.task(c.processLog)
		default:

		}
	}

}

// 将指定的日志文件打开，赋值给结构体
func (f *FileLogger) initFile() {
	logName := path.Join(f.filePath, f.fileName)
	// 打开文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("打开日志文件%s失败, %v", logName, err))
	}
	f.file = fileObj
	// 打开错误文件
	errLogFileName := fmt.Sprintf("%s.error", f.fileName)
	errLogFullFileName := path.Join(f.filePath, errLogFileName)
	errFileObj, err := os.OpenFile(errLogFullFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("打开错误日志文件%s失败, %v", errLogFullFileName, err))
	}
	f.errFile = errFileObj

}

// 检查是否要按文件大小拆分
func (f *FileLogger) checkSizeSplit(file *os.File) bool {

	if f.IsRotateBySize {
		info, _ := file.Stat()
		/*info, err := file.Stat()
		  if err != nil {
		  	return false
		  }*/
		fileSize := info.Size()

		return fileSize >= f.RotateBySize // 当传进来的日志文件大小超过maxSize，就返回true
	}
	return false
}

func (f *FileLogger) checkTimeSplit() bool {
	if f.IsRotateByTime {
		f.IsMakeRotateNumZero = f.CurrentTime() > recordTime
		return f.IsMakeRotateNumZero
	}
	return false
}

// 封装一个切分日志文件的方法
func (f *FileLogger) splitLogFile(ftype fileType, fileObj *os.File) {

	var (
		err       error
		rotateNum int
	)
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	//fmt.Println("111")

	switch ftype {
	case LOG_FILE:
		rotateNum = f.RotateLogNum
	case ERR_FILE:
		rotateNum = f.RotateErrNum
	}

	// 切分文件
	toCloseFileName := fileObj.Name() // 这里可以获取到文件的全路径
	toBackupFileName := fmt.Sprintf("%s.bak.%d", toCloseFileName, rotateNum)
	// 1. 把原来的文件关闭
	err = fileObj.Close()
	if err != nil {
		panic(fmt.Errorf("原有日志文件关闭失败:%s, %v", toCloseFileName, err))
	}
	// 2. 备份原来的文件
	err = os.Rename(toCloseFileName, toBackupFileName)
	//fmt.Printf("备份原来的文件: %s -> %s\n", toCloseFileName, toBackupFileName)
	if err != nil {
		panic(fmt.Errorf("归档日志文件失败:%s -> %s, %v", toCloseFileName, toBackupFileName, err))
	}
	// 3. 新建一个文件
	switch ftype {
	case LOG_FILE:
		f.fileName = generateLogFileName(f.filePrefix)
		toCloseFileName = path.Join(f.filePath, f.fileName)
	case ERR_FILE:
		toCloseFileName = path.Join(f.filePath, fmt.Sprintf("%s.error", f.fileName))
	}
	fileObj, err = os.OpenFile(toCloseFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("日志文件创建失败:%s, %v", toCloseFileName, err))
	}

	switch ftype {
	case LOG_FILE:
		f.file = fileObj
		f.RotateLogNum += 1
	case ERR_FILE:
		f.errFile = fileObj
		f.RotateErrNum += 1
	}

	recordTime = time.Now().Format(LOG_TIME_FORMAT)
	if f.IsMakeRotateNumZero {
		f.RotateLogNum = 0
		f.RotateErrNum = 0
	}
	//fmt.Println("222")
}

// 将公用的记录日志的功能封装成一个单独的方法
func (f *FileLogger) processLog(fileName string, line int, funcName string, level Level, format string, args ...interface{}) {
	if f.level > level {
		return
	}
	//f.file.Write()
	logFileBuf := bufio.NewWriter(f.file)
	errFileBuf := bufio.NewWriter(f.errFile)

	msg := fmt.Sprintf(format, args...) // 得到用户要记录的日志
	// 日志格式: [时间][文件:行号][函数名][日志级别]日志信息
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	//fileName, line, funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%d][%s][%s:%d][%s][%s]%s", Goid(), nowStr, fileName, line, funcName, logLevelStr, msg)
	// 往文件里写之前，要进行检查
	// 检查当前日志文件的大小是否超过了maxSize
	if f.checkSizeSplit(f.file) || f.checkTimeSplit() {
		f.splitLogFile(LOG_FILE, f.file)
	}
	fmt.Fprintln(logFileBuf, logMsg) // 利用fmt包将msg写入f.file文件中
	logFileBuf.Flush()
	// 如果是error或者fatal级别的日志，还要记录到f.errFile
	if level >= ErrorLevel {
		if f.checkSizeSplit(f.errFile) || f.checkTimeSplit() {
			f.splitLogFile(ERR_FILE, f.errFile)
		}
		fmt.Fprintln(errFileBuf, logMsg)
		errFileBuf.Flush()
	}
}

//func (f *FileLogger) watchFile() {
//	var (
//		fileName string
//		fileInfo os.FileInfo
//		err      error
//	)
//	var C = time.Tick(5 * time.Second)
//	for {
//		select {
//		case c := <-C:
//			fileName = f.file.Name()
//			fileInfo, err = os.Stat(fileName)
//			fmt.Printf("[%s] logFile %s size %d, err %v\n", c.Format(TIME_FULL_FORMAT), fileName, fileInfo.Size(), err)
//
//			fileName = f.errFile.Name()
//			fileInfo, err = os.Stat(fileName)
//			fmt.Printf("[%s] logFile %s size %d, err %v\n", c.Format(TIME_FULL_FORMAT), fileName, fileInfo.Size(), err)
//
//		default:
//		}
//	}
//}

func (f *FileLogger) Debug(msg string) {
	log := NewLog(DebugLevel, "%s", msg)
	f.jobChan <- log
}

func (f *FileLogger) Info(msg string) {
	log := NewLog(InfoLevel, "%s", msg)
	f.jobChan <- log
}

func (f *FileLogger) Warn(msg string) {
	log := NewLog(WarnLevel, "%s", msg)
	f.jobChan <- log
}

func (f *FileLogger) Error(msg string) {
	log := NewLog(ErrorLevel, "%s", msg)
	f.jobChan <- log
}

func (f *FileLogger) Fatal(msg string) {
	log := NewLog(FatalLevel, "%s", msg)
	f.jobChan <- log
}

// 方法 debug方法
func (f *FileLogger) Debugf(format string, args ...interface{}) {
	//f.log(DebugLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: DebugLevel, format: format, args: args}
	log := NewLog(DebugLevel, format, args...)
	f.jobChan <- log
}

// 方法 info方法
func (f *FileLogger) Infof(format string, args ...interface{}) {
	//f.log(InfoLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: InfoLevel, format: format, args: args}
	log := NewLog(InfoLevel, format, args...)
	f.jobChan <- log
}

// 方法 warn方法
func (f *FileLogger) Warnf(format string, args ...interface{}) {
	//f.log(WarnLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: WarnLevel, format: format, args: args}
	log := NewLog(WarnLevel, format, args...)
	f.jobChan <- log
}

// 方法 error方法
func (f *FileLogger) Errorf(format string, args ...interface{}) {
	//f.log(ErrorLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: ErrorLevel, format: format, args: args}
	log := NewLog(ErrorLevel, format, args...)
	f.jobChan <- log
}

// 方法 fatal方法
func (f *FileLogger) Fatalf(format string, args ...interface{}) {
	//f.log(FatalLevel, format, args...)
	//fileName, line, funcName := getCallerInfo(2)
	//log := Log{fileName: fileName, line: line, funcName: funcName, level: FatalLevel, format: format, args: args}
	log := NewLog(FatalLevel, format, args...)
	f.jobChan <- log
}

// Close()
func (f *FileLogger) Close() {
	close(f.jobChan)
	f.file.Close()
	f.errFile.Close()
}
