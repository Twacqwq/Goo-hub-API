// Package database 数据库操作
package database

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库 | gorm.Dialector 连接实体
func Connect(dbConfig gorm.Dialector, _logger gormlogger.Interface) {
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// 获取底层的sqlDB (通用的数据库接口)
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TableName(obj interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	stmt.Parse(obj)
	return stmt.Schema.Table
}
