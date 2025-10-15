package database

import (
	"fmt"
	"github.com/YuanJey/goutils2/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var Db *DataBase

type DataBase struct {
	dbFile    string
	conn      *gorm.DB
	viewMutex sync.RWMutex
	mRWMutex  sync.RWMutex
}

func init() {
	Db = &DataBase{dbFile: "./9zhou.db"}
	err := Db.initDB()
	if err != nil {
		panic(err)
	}
	fmt.Println("初始化成功")
}
func (d *DataBase) initDB() error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	db, err := gorm.Open(sqlite.Open(d.dbFile), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return utils.Wrap(err, "open db failed "+d.dbFile)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return utils.Wrap(err, "get sql db failed")
	}
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	d.conn = db
	db.AutoMigrate(
		&HxAccount{},
		&ShopAccount{},
		&Config{},
		&OrderStatistics{},
	)
	return nil
}
