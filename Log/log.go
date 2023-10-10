package Log

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
	"time"
)

type Logger struct {
	*logrus.Logger
}

func Initialize(path string, pattern string, options ...rotatelogs.Option) (error, ILogger) {
	var logInfo LoggerInfo = &DefaultLoggerInfo{}
	logInfo.SetName(pattern, nil)
	logInfo.SetPath(path)
	writer := logInfo.NewWriter(options...)

	//formater := logrus.JSONFormatter{
	//	FieldMap: logrus.FieldMap{
	//		logrus.FieldKeyLevel: "level",
	//		logrus.FieldKeyTime:  "timestamp",
	//	},
	//}
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "create_time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function",
		},
	}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		return frame.Function, frame.File
	}
	formatter.TimestampFormat = "2006-01-02 15:04:05 Z07:00"

	log := &Logger{logrus.New()}
	log.Formatter = formatter
	log.Out = writer
	log.SetReportCaller(true)

	return nil, log
}

func (l *Logger) GetOutput() io.Writer {
	return l.Out
}

type ILogger interface {
	GetOutput() io.Writer
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	TraceFn(fn logrus.LogFunction)
	DebugFn(fn logrus.LogFunction)
	InfoFn(fn logrus.LogFunction)
	PrintFn(fn logrus.LogFunction)
	WarnFn(fn logrus.LogFunction)
	WarningFn(fn logrus.LogFunction)
	ErrorFn(fn logrus.LogFunction)
	FatalFn(fn logrus.LogFunction)
	PanicFn(fn logrus.LogFunction)

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type IEntry interface {
	WithError(err error) IEntry
	WithContext(ctx context.Context) IEntry
	WithField(key string, value interface{}) IEntry
	WithFields(fields logrus.Fields) IEntry
	WithTime(t time.Time) IEntry

	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	sprintlnn(args ...interface{}) string
}
