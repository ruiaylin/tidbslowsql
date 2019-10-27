package components

import (
	"log"

	"github.com/BurntSushi/toml"
)

func InitConfig() {
	path := "conf/config.ini"
	//if wd, err := os.Getwd(); err == nil && strings.Index(wd, "pano/") > -1 {
	//
	//	path = wd[:strings.Index(wd, "pano/")+6] + path
	//}
	_, err := toml.DecodeFile(path, &Config)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	//	Config.reflect()
}

// Config config
var Config config

type config struct {
	DB struct {
		MysqlServerRead  string
		MysqlServerWrite string
		LogFlag          bool
	}

	Log struct {
		Remote     bool
		RemoteAddr string
		Level      int
	}
}
