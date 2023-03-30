# android Project 2 后端

## 文档

[安卓实验二答辩报告文档](https://cmjhgnav4v.feishu.cn/docx/TSv9dfxJpoXgU6xzi4ucUnprnbh)

## 使用
### 配置数据库
自行在 config 文件夹下创建本地数据库配置文件 conf_db.go ,作为配置本地MySql和Redis文件
```
| - db
    | - config.go
```
```go
package config

// MySql
// 需要修改成自己的配置
var (
	user     = "111"
	password = "123456"
	db_name  = "androidProject2"
)

// Redis
// 需要修改成自己的配置
var (
	RedisIP       = "localhost"
	RedisPort     = "6379"
	RedisAddr     = RedisIP + ":" + RedisPort
	RedisPassword = ""
	RedisDB       = 0
)
```


## 运行
```go
run main.go
```