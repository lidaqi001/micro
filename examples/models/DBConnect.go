package models

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() *gorm.DB {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "sxx_admin:sxx170615@tcp(rm-bp196c94547xbb474lo.mysql.rds.aliyuncs." +
		"com:3306)/shanxiuxia3_product?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql connect failed : " + err.Error())
		logger.Fatal("mysql connect failed : " + err.Error())
	}
	return db
}
