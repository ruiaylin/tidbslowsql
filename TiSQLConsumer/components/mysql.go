package components

import (
	"sync"
        "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// DBOnce control for the initialize
	DBOnce sync.Once
	// DBReader Reader
	DBReader *gorm.DB
	// DBWriter writer
	DBWriter *gorm.DB
)

// InitDB init the db connection
func InitDB() {
	DBOnce.Do(func() {
		var err error
		DBReader, err = gorm.Open("mysql", Config.DB.MysqlServerRead)
		if err != nil {
			fmt.Printf("Init DB error : %v", err)
			return
		}
		DBReader.LogMode(Config.DB.LogFlag)
		DBWriter, err = gorm.Open("mysql", Config.DB.MysqlServerWrite)
		if err != nil {
			fmt.Printf("Init DB error : %v", err)
			return
		}
		DBReader.LogMode(Config.DB.LogFlag)
		DBWriter.LogMode(Config.DB.LogFlag)
	})

}
