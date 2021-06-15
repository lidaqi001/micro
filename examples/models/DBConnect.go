package models

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() *gorm.DB {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.M_USER,
		config.M_PASSWORD,
		config.M_HOST,
		config.M_DBNAME,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql connect failed : " + err.Error())
		logger.Fatal("mysql connect failed : " + err.Error())
	}
	return db
}
