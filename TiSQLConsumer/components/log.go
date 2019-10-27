package components

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"path"
	"runtime"
	"sync"
)

var Logger *logrus.Logger
var writer io.Writer
var once sync.Once

//初始化日志
func InitLogger() {
	once.Do(func() {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		if Config.Log.Remote {
			writer = &udpWriter{}
		} else {
			writer = os.Stderr
		}
		Logger = logrus.New()
		hook := []logrus.Hook{&ContextHook{}}
		Logger.Hooks = logrus.LevelHooks{
			logrus.PanicLevel: hook,
			logrus.FatalLevel: hook,
			logrus.ErrorLevel: hook,
			logrus.WarnLevel:  hook,
			logrus.InfoLevel:  hook,
			logrus.DebugLevel: hook,
		}
		logrus.SetOutput(writer)
		logrus.SetLevel(logrus.Level(Config.Log.Level))

		Logger.Info("日志初始化完毕")
	})
}

func GetWriter() io.Writer {
	return writer
}

//Log HOOKS
type ContextHook struct{}


func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(10); ok {
		funcName := runtime.FuncForPC(pc).Name()

		entry.Data["source"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(funcName))
	}

	return nil
}

//udpWriter实现了Write接口
//日志如果为远程日志，则调用此结构
//log输出会自动调用Write方法
//TODO 是否需要队列输出？
type udpWriter struct{}

func (udp *udpWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	dial, err := net.Dial("udp", Config.Log.RemoteAddr)
	if err != nil {
		return -1, err
	}
	defer dial.Close()
	return dial.Write(p)
}
