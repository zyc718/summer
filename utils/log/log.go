package log

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"summer/utils/log/writer"
)

type logConfig struct {
	Enable    bool                       `json:"enable"`
	Level     string                     `json:"level"`
	UseWriter string                     `json:"use_writer"`
	Writers   map[string]json.RawMessage `json:"writers"`
}

func Init(configRaw json.RawMessage) error {
	//解析json配置文件
	var config logConfig
	err := json.Unmarshal(configRaw, &config)
	if err != nil {
		return err
	}

	//进行日志初始化
	//进行日志前 hook 拦截 设置文件、目录地址等
	logrus.AddHook(NewContextHook())

	//设置日志内容使用json格式
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 输出到stdout标准输出， 参数可以使任何一个io.Writer实现
	writers := []io.Writer{}

	if config.Enable {
		var logWriter io.Writer
		fmt.Printf("qqqq %v\n", config.UseWriter)
		switch config.UseWriter {
		case "file":
			logWriter, err = writer.NewFileWriter(config.Writers[config.UseWriter])
			if err != nil {
				return err
			}
		}
		writers = append(writers, logWriter)
	}
	writers = append(writers, os.Stdout)
	multiWriter := io.MultiWriter(writers...)
	logrus.SetOutput(multiWriter)

	l := logrus.GetLevel()

	switch config.Level {
	case "panic":
		l = logrus.PanicLevel
	case "fatal":
		l = logrus.FatalLevel
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	case "trace":
		l = logrus.TraceLevel

	}
	logrus.SetLevel(l)
	return nil
}

func Tracef(format string, v ...interface{}) {
	logrus.Tracef(format, v...)
}
func Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logrus.Fatalf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	logrus.Panicf(format, v...)
}
func CTracef(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Tracef(format, v...)
}

func CDebugf(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Debugf(format, v...)
}

func CInfof(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Infof(format, v...)
}

func CWarnf(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Warnf(format, v...)
}

func CErrorf(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Errorf(format, v...)
}

func CFatalf(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Fatalf(format, v...)
}

func CPanicf(code int, format string, v ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"code": strconv.Itoa(code),
	}).Panicf(format, v...)
}
