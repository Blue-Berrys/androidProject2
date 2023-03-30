package model

import (
	"androidProject2/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlUrl(user string, pwd string, db_name string) string {
	return user + ":" + pwd + "@tcp(localhost:3306)/" + db_name + "?charset=utf8&parseTime=True&loc=Local"
}

var DSN = mysqlUrl(config.SqlUser, config.SqlPassword, config.SqlDb_name)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{
		SkipDefaultTransaction: true, //关闭默认事务
		PrepareStmt:            true, //缓存预编译语句
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&User{}, &FriendsChat{}, &Like{}, &Comment{})
	if err != nil {
		panic(err)
	}
}
