package main

import (
	"fmt"
	"github.com/lidaqi001/micro/examples/models"
	"gorm.io/gorm"
)

// 这里仅做简单示例
// 更多用法，查看GORM文档 https://gorm.io/zh_CN/docs/

// 模型声明
type Order struct {
	Id         uint32
	Number     string
	Customer   string
	CustomerId uint32
	Cellphone  string
}

// gorm默认会为你设置的表模型名称后加一个s
// 例如 ：上面设置的 Order =》 查询时会变成 Orders
// 所以手动指定表名

// 当然还有第二种方法，查询时指定表名：
// db.Table("order")
func (Order) TableName() string {
	return "order"
}

var db *gorm.DB

func main() {

	/*******************************************************
				            获取DB连接
	*******************************************************/
	db = models.DB()

	//add := create()
	//order := find(add.Id)
	//update(order)
	//delete(order)

}

/*******************************************************
			            Create
*******************************************************/
func create() Order {
	order := Order{
		Number:     "test1111",
		Customer:   "liqi",
		CustomerId: 99999,
		Cellphone:  "18888888888",
	}

	addResult := db.Create(&order) // 通过数据的指针来创建
	fmt.Println("返回插入数据的主键：", order.Id)
	fmt.Println("create返回error：", addResult.Error)
	fmt.Println("插入记录的条数：", addResult.RowsAffected)
	fmt.Println("插入数据：", order)
	return order
}

/*******************************************************
			            Find
*******************************************************/
func find(id uint32) Order {
	// 单条查询
	var order Order
	_ = db.Where("id = ?", id).First(&order)
	fmt.Println("data : ", order)

	// 多条查询
	var orders []Order
	_ = db.Where("id <= ?", 10).Find(&orders)
	fmt.Println("datas : ", orders)

	return order
}

/*******************************************************
			            Update
*******************************************************/
func update(order Order) {

	order.Customer = "lidaqi"
	updateRes := db.Save(&order)

	fmt.Println("update返回error：", updateRes.Error)
	fmt.Println("更新记录的条数：", updateRes.RowsAffected)
}

/*******************************************************
			            Delete
*******************************************************/
func delete(order Order) {

	deleteRst := db.Delete(&order)
	fmt.Println("delete返回error：", deleteRst.Error)
	fmt.Println("删除记录的条数：", deleteRst.RowsAffected)
}
