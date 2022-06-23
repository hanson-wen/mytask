package dao

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// init connect db
func init() {
	var err error
	user := beego.AppConfig.String("mysql_user")
	pass := beego.AppConfig.String("mysql_pass")
	host := beego.AppConfig.String("mysql_host")
	port := beego.AppConfig.String("mysql_port")
	dbName := beego.AppConfig.String("mysql_db")

	links := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	logs.Info("db init:", links)
	db, err = gorm.Open("mysql", links)
	if err != nil {
		logs.Error("db init err:", err.Error())
	} else {
		logs.Error("db init success")
	}
	db.LogMode(true)
	//defer db.Close()
}

// CreateData common create function
func CreateData(m interface{}) (err error) {
	err = db.Model(m).Create(m).Error
	if err != nil {
		logs.Error("create data err:%v, %s", m, err.Error())
	}
	return
}

// ModifyData common modify data, need primary
func ModifyData(m interface{}) (err error) {
	err = db.Model(m).Update(m).Error
	if err != nil {
		logs.Error("modify data err:%v, %s", m, err.Error())
	}
	return
}
