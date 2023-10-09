package Log

import (
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
)

type LoggerInfo interface {
	SetName(string, LoggerNameHandler)
	GetName() string
	SetPath(string)
	GetSuffix() string
	GetFullPath() string
	NewWriter(...rotatelogs.Option) io.Writer
}

type LoggerNameHandler interface {
	Handle(LoggerInfo) string
}

type LoggerNameExtension func(LoggerInfo) string

func (lne LoggerNameExtension) Handle(li LoggerInfo) string {
	return lne(li)
}

var defaultLnx = LoggerNameExtension(func(li LoggerInfo) string {

	return li.GetName() + "." + li.GetSuffix()
})

type DefaultLoggerInfo struct {
	path   string
	name   string
	suffix string
	nex    LoggerNameHandler
}

func (dli *DefaultLoggerInfo) SetName(pattern string, f LoggerNameHandler) {

	fnBytes := []byte(pattern)
	n := bytes.LastIndexByte(fnBytes, byte('.'))
	dli.suffix = string(fnBytes[n+1:])
	dli.name = string(fnBytes[:n])

	dli.nex = defaultLnx
	if f != nil {
		dli.nex = f
	}
}

func (dli *DefaultLoggerInfo) SetPath(path string) {
	dli.path = path
}

func (dli *DefaultLoggerInfo) GetName() string {
	return dli.name
}

func (dli *DefaultLoggerInfo) GetSuffix() string {
	return dli.suffix
}

func (dli *DefaultLoggerInfo) GetFullPath() string {
	if dli.nex != nil {
		return dli.path + "/" + dli.nex.Handle(dli)
	}
	return ""
}

func (dli *DefaultLoggerInfo) NewWriter(options ...rotatelogs.Option) io.Writer {
	if _, err := os.Stat(dli.path); os.IsNotExist(err) {
		err = os.MkdirAll(dli.path, 0755)
		if err != nil {
			fmt.Println(err)
			return os.Stderr
		}
	}
	filename := dli.GetFullPath()
	//fmt.Println(filename)
	//writer, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	writer, err := rotatelogs.New(filename, options...)

	if err != nil {
		fmt.Println(err)
		return os.Stderr
	}
	return writer
}

//type ILoggerWriter interface {
//	Write(p []byte) (int, error)
//	Close() error
//}
//
//type LoggerWriter struct {
//	path  string
//	fw    *os.File
//	ln    LoggerName
//	fn    LoggerNameHandler
//	mutex sync.Mutex
//}
//
//func (lw *LoggerWriter) SetName(pattern string) {
//	ln := &DefaultLoggerName{}
//	ln.SetName(pattern)
//	lw.ln = ln
//	lw.fn = defaultNameExt
//}
//
//func (lw *LoggerWriter) SetNameExt(pattern string, f LoggerNameHandler) {
//	ln := &DefaultLoggerName{}
//	ln.SetName(pattern)
//	lw.ln = ln
//	lw.fn = f
//}
//
//func (lw *LoggerWriter) SetPath(path string) {
//	lw.path = path
//}
//
//func (lw *LoggerWriter) Write(p []byte) (int, error) {
//	lw.mutex.Lock()
//	defer lw.mutex.Unlock()
//	defer lw.Close()
//	writer := lw.NewWriter()
//	fmt.Println("ssss")
//	return writer.Write(p)
//}
//
//func (lw *LoggerWriter) Close() error {
//	if lw.fw == nil{
//		return errors.New("writer is nill")
//	}
//	fmt.Println("close")
//	return lw.fw.Close()
//}
//
//func (lw *LoggerWriter) NewWriter() io.Writer {
//	if _, err := os.Stat(lw.path); os.IsNotExist(err) {
//		err = os.MkdirAll(lw.path, 0755)
//		if err != nil {
//			fmt.Println(err)
//			return os.Stderr
//		}
//	}
//	filename := lw.path + "/" + lw.fn.Handle(lw.ln)
//	//fmt.Println(filename)
//	//writer, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
//	writer, err := rotatelogs.New(filename, rotatelogs.WithClock(rotatelogs.Local))
//
//	if err != nil {
//		fmt.Println(err)
//		return os.Stderr
//	}
//	return writer
//}
