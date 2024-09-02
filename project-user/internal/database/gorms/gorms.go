package gorms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"ms_project/project-user/internal/model"
	"time"
)

var _db *gorm.DB

func Database(connString string) {
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		fmt.Println(err)
		panic("mysql数据库连接错误")
	}
	fmt.Println("mysql数据库连接成功")
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	db.SingularTable(true)       //创建的表名不加s
	db.DB().SetMaxIdleConns(20)  //设置连接池
	db.DB().SetMaxOpenConns(100) //设置最大连接数
	db.DB().SetConnMaxLifetime(time.Second * 30)
	_db = db
	//migration()
}
func migration() {
	//自动迁移模式
	_db.Set("gorm:table_options", "charset=utf8mb4,ENGINE=InnoDB").
		AutoMigrate(&model.Member{},
			&model.Organization{},
		)
}
func NewDBClient() *gorm.DB {
	DB := _db
	return DB
}
