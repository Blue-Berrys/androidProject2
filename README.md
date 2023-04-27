# android Project 2 后端

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | 安卓课程设计大作业2 后端，使用 Gin 和 Gorm 构建，由 [Blue-Berrys](https://github.com/Blue-Berrys) 开发。 |
| -------------------------------------------------------------------- |------------------------------------------------------------------------------------|
## 详细说明文档

[安卓实验二答辩报告文档](https://cmjhgnav4v.feishu.cn/docx/TSv9dfxJpoXgU6xzi4ucUnprnbh)

## 项目结构
* [cache](https://github.com/Blue-Berrys/androidProject2/tree/main/cache) - 缓存
  * [redis](https://github.com/Blue-Berrys/androidProject2/tree/main/cache/Redis) - 数据库缓存
  * [minio](https://github.com/Blue-Berrys/androidProject2/tree/main/cache/minio) - 对象存储服务OSS
* [config](https://github.com/Blue-Berrys/androidProject2/tree/main/config) - 配置文件
  * [conf.go](https://github.com/Blue-Berrys/androidProject2/blob/main/config/conf.go) - 常用返回码
  * conf_db.go - 需要自己配置，Mysql、Redis、Minio服务器地址，[在下面有介绍](#配置数据库)
* [handler](https://github.com/Blue-Berrys/androidProject2/tree/main/handler) - 解析请求参数
  * [Comment](https://github.com/Blue-Berrys/androidProject2/tree/main/handler/Comment) - 与评论相关的handler
    * [comment_action_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/Comment/comment_action_handler.go) - 增加或删除评论的handler
    * [comment_list_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/Comment/comment_list_handler.go) - 评论列表的handler
  * [FriendsChat](https://github.com/Blue-Berrys/androidProject2/tree/main/handler/FriendsChat) - 与朋友圈相关的handler
    * [publish_action_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/FriendsChat/publish_action_handler.go) - 发布朋友圈的handler
    * [publish_list_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/FriendsChat/publish_list_handler.go) - 朋友圈列表的handler
    * [delete_action_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/FriendsChat/delete_action_handler.go) - 删除朋友圈的handler
  * [like](https://github.com/Blue-Berrys/androidProject2/tree/main/handler/like) - 与点赞相关的handler
    * [like_action_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/like/like_action_handler.go) - 点赞/取消赞操作的handler
    * [like_list_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/like/like_list_handler.go) - 点赞列表的handler
  * [user](https://github.com/Blue-Berrys/androidProject2/tree/main/handler/user) 与用户相关的handler
    * [level_action_handler](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/user/level_action_handler.go) - 修改用户等级的handler
    * [level_list_handler](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/user/level_list_handler.go) - 列出不同等级所有用户的handler
    * [user_info_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/user/user_info_handler.go) - 用户信息的handler
    * [user_login_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/user/user_login_handler.go) - 用户登录的handler
    * [user_register_handler.go](https://github.com/Blue-Berrys/androidProject2/blob/main/handler/user/user_regitser_handler.go) - 用户注册的handler
* [middleware](https://github.com/Blue-Berrys/androidProject2/tree/main/middleware) - 中间件
  * [Bcrypt](https://github.com/Blue-Berrys/androidProject2/tree/main/middleware/Bcrypt) - 加盐加密
    * [bcrypt.go](https://github.com/Blue-Berrys/androidProject2/blob/main/middleware/Bcrypt/bcrypt.go) - 对明文密码加盐加密
    * [bcrypt_test.go](https://github.com/Blue-Berrys/androidProject2/blob/main/middleware/Bcrypt/bcrypt_test.go) - 测试加密正确性
  * [JWT](https://github.com/Blue-Berrys/androidProject2/tree/main/middleware/JWT) - 用户鉴权token
    * [JWT.go](https://github.com/Blue-Berrys/androidProject2/blob/main/middleware/JWT/JWT.go) - 生成用户鉴权token
    * [JWT_test.go](https://github.com/Blue-Berrys/androidProject2/blob/main/middleware/JWT/JWT_test.go) - 测试生成用户鉴权token正确性
* [model](https://github.com/Blue-Berrys/androidProject2/tree/main/model) - 数据库
  * [db](https://github.com/Blue-Berrys/androidProject2/tree/main/model/db) 表
    * [initTable.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/db/initTable.go) - 初始化表
    * [table.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/db/table.go) - 设计数据库表
    * [initTable_test.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/db/initTable_test.go) - 测试表是否建立正确
  * [Comment](https://github.com/Blue-Berrys/androidProject2/tree/main/model/Comment) - 与评论有关的数据库操作
    * [comment.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/Comment/comment.go) - 操作评论数据库函数
  * [friendschat](https://github.com/Blue-Berrys/androidProject2/tree/main/model/friendschat) - 与朋友圈有关的数据库操作
    * [friendschat.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/friendschat/friendschat.go) - 操作朋友圈数据库函数
  * [like](https://github.com/Blue-Berrys/androidProject2/tree/main/model/like) - 与点赞有关的数据库操作
    * [like.go](https://github.com/Blue-Berrys/androidProject2/blob/main/model/like/like.go) - 操作点赞数据库函数
  * [user](https://github.com/Blue-Berrys/androidProject2/tree/main/model/user) - 与用户有关的数据库操作
    * [user](https://github.com/Blue-Berrys/androidProject2/blob/main/model/user/user.go) - 操作用户数据库函数
* [router](https://github.com/Blue-Berrys/androidProject2/tree/main/router) - 路由
  * [router.go](https://github.com/Blue-Berrys/androidProject2/blob/main/router/router.go) - 自定义路由接口
* [service](https://github.com/Blue-Berrys/androidProject2/tree/main/service) - 参数合法性校验，调用数据库
  * [Comment](https://github.com/Blue-Berrys/androidProject2/tree/main/service/Comment) - 与评论有关的service
    * [comment_action_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/Comment/comment_action_service.go) - 增加或删除评论的service
    * [comment_list_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/Comment/comment_list_service.go) - 评论列表的service
  * [FriendsChat](https://github.com/Blue-Berrys/androidProject2/tree/main/service/FriendsChat) - 与朋友圈有关的service
    * [publish_action_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/FriendsChat/publish_action_service.go) - 发布朋友圈的service
    * [publish_list_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/FriendsChat/publish_list_service.go) - 朋友圈列表的service
    * [delete_action_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/FriendsChat/delete_action_service.go) - 删除朋友圈的service
  * [Like](https://github.com/Blue-Berrys/androidProject2/tree/main/service/Like) - 与点赞有关的service
    * [like_action_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/Like/like_action_service.go) - 点赞/取消赞操作的service
    * [like_list_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/Like/like_list_service.go) - 点赞列表的service
  * [user](https://github.com/Blue-Berrys/androidProject2/tree/main/service/user) - 与用户有关的service
    * [level_action_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/user/level_action_service.go) - 用户等级切换的service
    * [level_list_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/user/level_list_service.go) - 用户不同等级列表的service
    * [user_info_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/user/user_info_service.go) - 用户信息的service
    * [user_login_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/user/user_login_service.go) - 用户登录的service
    * [user_register_service.go](https://github.com/Blue-Berrys/androidProject2/blob/main/service/user/user_register_service.go) - 用户注册的service
* [util](https://github.com/Blue-Berrys/androidProject2/tree/main/util) - 套Json结构
  * [table.go](https://github.com/Blue-Berrys/androidProject2/blob/main/util/table.go) - 返回接口需要的Json结构
* [main.go](https://github.com/Blue-Berrys/androidProject2/blob/main/main.go) - 主函数

## 准备环境
* Linux / MacOs / Windows
* Go v1.9
* MySql v8.0.31
* Redis v7.0.7
* Minio 

## 使用
### 配置数据库
自行在 config 文件夹下创建本地数据库配置文件 conf_db.go,需要本地配置好MySql、Redis、Minio
```
| - config
    | - conf_db.go
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

// minio
var (
	Miniourl       = "10.215.221.49:9000" //搭建容器地址
	MinioaccessKey = "minioadmin"         //minioadmin
	MiniosecretKey = "minioadmin"         //minioadmin //key
	HeartbeatTime  = 2 * 60

	BucketName = "images"
	Location   = "cn-north-1"
)

// PlayUrlPrefix 存储的图片的链接
var PlayUrlPrefix = "/images/"
```


## 运行
``` go
go run main.go
```

## 压力测试结果
每个接口各有3、4组测试样例，均测试500次，总共 HTTP 接口请求数 15500，总耗时 1626.7秒，平均接口请求耗时 94ms ，QPS 大约 100