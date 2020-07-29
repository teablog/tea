package db

import (
	"github.com/teablog/tea/internal/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

var DB *gorm.DB
var DB_ONCE sync.Once

func NewDB(dsn string){
	DB_ONCE.Do(func() {
		// videos_t
		var err error
		DB, err = gorm.Open("mysql", dsn)
		if err != nil {
			logger.Fatal("gorm open %s", err)
		}
		if err = DB.DB().Ping(); err != nil {
			logger.Fatal(err)
		}
		DB.SingularTable(true)
	})
}

func Close()  {
	_ = DB.Close()
}
