package models

import (
	"ctf/config"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	CreatedAt MyTime  `gorm:"column:created_at;type:datetime;" json:"created_at"`
	UpdatedAt MyTime  `gorm:"column:updated_at;type:datetime;" json:"updated_at"`
	DeletedAt *MyTime `gorm:"column:deleted_at;type:datetime;" json:"deleted_at" sql:"index"`
}
type MyTime struct {
	time.Time
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t MyTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *MyTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = MyTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

const (
	sqlType = "mysql"
)

func Setup() {
	var err error
	//db, err = gorm.Open(sqlType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
	//	sqlUser,
	//	sqlPassword,
	//	sqlHost,
	//	dbName))
	db, err = gorm.Open(sqlType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Global.MysqlInfo.User,
		config.Global.MysqlInfo.Password,
		config.Global.MysqlInfo.Host,
		config.Global.MysqlInfo.Db))

	if err != nil {
		log.Printf(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Global.MysqlInfo.User,
			config.Global.MysqlInfo.Password,
			config.Global.MysqlInfo.Host,
			config.Global.MysqlInfo.Db))
		log.Fatalf("models.Setup  err: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	db.SingularTable(true)
	//db.Callback().Create().Replace("gorm:insert_option", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_option", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	//db.LogMode(true)
	//创建表
	db.AutoMigrate(
		BlockScan{},
		UserTable{},
		UserLowersTable{},
	)
	//CreateTableByHash()
}
