package log_utils

import (
	"fmt"
	"github.com/tjlcast/go_common/rand_utils"
	"path"
	"runtime"
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
}

func ShowReportCaller() {
	Logger.SetReportCaller(true)
}

func SetLogLevel(level log.Level) {
	Logger.SetLevel(level)
}

func getTimestampStr() string {
	now := time.Now()
	dateString := fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return dateString
}

func TjlTestLog(msg string) {
	// This is for tjl's test infomation, Will remove in future.
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	msg = fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tjl tes: %s", msg)
	msg = fmt.Sprintf("\n%s [TJL]\t :%s  at (%s:%d [Method %s])\n", getTimestampStr(), msg, file, line, f.Name())
	fmt.Println(Yellow(msg))
}

func TestLog(msg string) {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	msg = fmt.Sprintf("\n%s [Test]\t :%s  at (%s:%d [Method %s])\n", getTimestampStr(), msg, file, line, f.Name())
	fmt.Println(Yellow(msg))
}

func ModuleLog(mod string, msg string) {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	msg = fmt.Sprintf("\n%s [%s]\t :%s  at (%s:%d [Method %s])\n", getTimestampStr(), mod, msg, file, line, f.Name())
	fmt.Println(Yellow(msg))
}

func RandLog(msg string, r int) {
	if r < 0 {
		r = 10
	}
	if int(rand_utils.GenRandInt(r)) == 0 {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	msg = fmt.Sprintf("\n%s [Rand%dLog]\t :%s  at (%s:%d [Method %s])\n", getTimestampStr(), r, msg, file, line, f.Name())
	fmt.Println(Blue(msg))
}

func Wrap(domain string, format string, a ...interface{}) string {
	return wrap(domain, format, a...)
}

func wrap(domain string, format string, a ...interface{}) string {
	format = "["+domain+"] " + format
	if a==nil || len(a)==0 {
		return format
	}
	return fmt.Sprintf(format, a...)
}

func TestLogWrap(format string, a ...interface{}) string {
	return wrap("Test", format, a...)
}

func MasterLogWrap(format string, a ...interface{}) string {
	return wrap("Master", format, a...)
}

func WorkerLogWrap(format string, a ...interface{}) string {
	return wrap("Worker", format, a...)
}

func BootLogWrap(format string, a ...interface{}) string {
	return wrap("Boot", format, a...)
}

func RepoLogWrap(format string, a ...interface{}) string {
	return wrap("Repo", format, a...)
}

func ConsisLogWrap(format string, a ...interface{}) string {
	return wrap("Consis", format, a...)
}

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}