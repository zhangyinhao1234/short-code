package initialize

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"short-code/global"
	"time"
)

var (
	env       = os.Getenv("k8s_env")
	podName   = os.Getenv("podName")
	namespace = os.Getenv("namespace")
)

func Logrus() {
	logPath := global.CONF.Logger.FileDir   // 日志存放路径
	linkName := global.CONF.Logger.LinkName // 最新日志的软连接路径
	level, _ := logrus.ParseLevel(global.CONF.Logger.Level)
	global.Logger.SetLevel(level) // 设置日志级别
	logWriter, _ := rotatelogs.New(
		logPath+"%Y%m%d.log",                      // 日志文件名格式
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 最多保留7天之内的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 一天保存一个日志文件
		rotatelogs.WithLinkName(linkName),         // 为最新日志建立软连接
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 格式日志时间
	})
	global.Logger.AddHook(Hook)
	if global.CONF.Server.StartEnv == "local" {
		global.LOG = logrus.NewEntry(global.Logger)
	} else {
		global.LOG = logrus.NewEntry(global.Logger).WithFields(logrus.Fields{"appName": global.CONF.Server.Name,
			"podName":      podName,
			"podNamespace": namespace, "env": env})
	}

}
