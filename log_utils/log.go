package log_utils

import (
	"path"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func SetLog(maxAge time.Duration, duration time.Duration) {
	ConfigLocalFilesystemLogger("log", "daily", maxAge, duration, Logger)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration, Logger *log.Logger) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+"-%Y%m%d%H%M.log",
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		Logger.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	errWriter, err := rotatelogs.New(
		baseLogPath+"-error-%Y%m%d%H%M.log",
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		Logger.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: errWriter,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{ForceColors: true})

	Logger.SetReportCaller(true) //将函数名和行数放在日志里面
	Logger.AddHook(lfHook)

	Logger.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	//Logger.SetLevel(log.DebugLevel)
	Logger.SetLevel(log.InfoLevel) // default runtime.
}

func init() {

	SetLog(time.Hour*24*7, time.Hour*24)
	Logger.Info("Config Logger Done")
}

func ShowReportCaller() {
	Logger.SetReportCaller(true)
}

func SetLogLevel(level log.Level) {
	Logger.SetLevel(level)
}

func TjlTestLog(msg string) {
	Logger.Info(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tjl test: " + msg)
}
