package main

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func LogSetUp() {
	defaultLogger = logrus.New()
	defaultLogger.SetLevel(logrus.DebugLevel)

	_ = SetupDefaultLogger("log", "generate", "debug",
		604800, 86400)
}

// SetupDefaultLogger configure logger instance with user provided settings
func SetupDefaultLogger(logPath string, fileName string, level string, maxAge int, rotationTime int) (err error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	defaultLogger.SetLevel(lvl)

	ConfigLocalFilesystemLogger(defaultLogger, logPath, fileName, maxAge, rotationTime)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func ConfigLocalFilesystemLogger(l *logrus.Logger, logPath string, logFileName string, maxAge int, rotationTime int) {

	maxAgeDuration := time.Second * time.Duration(int64(maxAge))
	rotationTimeDuration := time.Second * time.Duration(int64(rotationTime))

	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d",
		//rotatelogs.WithLinkName(baseLogPaht),    			// 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAgeDuration),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTimeDuration), // 日志切割时间间隔
	)
	if err != nil {
		l.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})

	//注意上面这个 and &amp;符号被转义了
	l.AddHook(lfHook)
}

// Logging struct that holds all user configurable options for the logger
type Logging struct {
	Enabled              bool   `json:"enabled,omitempty"`
	File                 string `json:"file"`
	ColourOutput         bool   `json:"colour"`
	ColourOutputOverride bool   `json:"colourOverride,omitempty"`
	Level                string `json:"level"`
	Rotate               bool   `json:"rotate"`
	LogPath              string `json:"logpath"`
	MaxAge               int    `json:"maxage"`
	Rotationtime         int    `json:"rotationtime"`
}

var (
	defaultLogger *logrus.Logger
	//// LogPath location to store logs in
	LogPath string
	//// 默认参数；可以在外部重载参数
	Logger = &Logging{
		Enabled:              true,
		File:                 "default.log",
		ColourOutput:         false,
		ColourOutputOverride: false,
		Level:                "DEBUG",
		Rotate:               true,
		LogPath:              "./log",
		MaxAge:               604800000000000,
		Rotationtime:         86400000000000,
	}
)

// Info handler takes any input returns unformatted output to infoLogger writer
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infof handler takes any input infoLogger returns formatted output to infoLogger writer
func Infof(data string, v ...interface{}) {
	defaultLogger.Infof(data, v...)
}

// Infoln handler takes any input infoLogger returns formatted output to infoLogger writer
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Print aliased to Standard log.Print
var Print = defaultLogger.Print

// Printf aliased to Standard log.Printf
var Printf = defaultLogger.Printf

// Println aliased to Standard log.Println
var Println = defaultLogger.Println

// Debug handler takes any input returns unformatted output to infoLogger writer
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Debugf handler takes any input infoLogger returns formatted output to infoLogger writer
func Debugf(data string, v ...interface{}) {
	defaultLogger.Debugf(data, v...)
}

// Debugln handler takes any input infoLogger returns formatted output to infoLogger writer
func Debugln(v ...interface{}) {
	defaultLogger.Debugln(v...)
}

// Warn handler takes any input returns unformatted output to warnLogger writer
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warnf handler takes any input returns unformatted output to warnLogger writer
func Warnf(data string, v ...interface{}) {
	defaultLogger.Warnf(data, v...)
}

// Error handler takes any input returns unformatted output to errorLogger writer
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorf handler takes any input returns unformatted output to errorLogger writer
func Errorf(data string, v ...interface{}) {
	defaultLogger.Errorf(data, v...)
}

// Fatal handler takes any input returns unformatted output to fatalLogger writer
func Fatal(v ...interface{}) {
	// Send to Output instead of Fatal to allow us to increase the output depth by 1 to make sure the correct file is displayed
	defaultLogger.Fatal(v...)
	os.Exit(1)
}

// Fatalf handler takes any input returns unformatted output to fatalLogger writer
func Fatalf(data string, v ...interface{}) {
	// Send to Output instead of Fatal to allow us to increase the output depth by 1 to make sure the correct file is displayed
	defaultLogger.Fatalf(data, v...)
	os.Exit(1)
}
