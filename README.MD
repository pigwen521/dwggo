## 包含功能模块
- mvc
- http微框架：gin github.com/gin-gonic/gin
- orm：gorm.io/gorm （表名不为复数，SingularTable）
- 日志：zap lumberjack
- mysql,redis 连接池
- 配置文件 github.com/spf13/viper

## 功能介绍
### MVC
请求./ademo/query?name=test，对应业务代码如下：
- 控制器 ./controller/AdemoController.go
- 数据输入输出模型 ./model/AdemoModel.go
- 业务逻辑 ./service/AdemoService.go
- 数据实体 ./entity/Ademo.go
### 配置文件
- 默认读取./config/config.yaml secret.yaml system.yaml3个文件
- 使用方法：core.GetConfigString("app.name") core.GetConfigInt("app.port")
### 日志
保存目录、文件名、分割大小见配置文件logger项，包含：
- error,info,debug日志
- db日志 
- 启动日志
- 请求API日志
#### 使用方法
- core.LogError("log...")
- core.LogInfo("log...")
- core.LogDebug("log...")
- core.LogErrorAndPanic("log...")
#### 自定义日志文件
- core.LogErrorCustom("log...",FILE_NAME)
- core.LogInfoCustom("log...",FILE_NAME)
### 限流器
- 内置的limiter，令牌生成速率和容量见配置文件limiter
### API接口请求
默认会记录日志（请求参数和响应返回内容）
- GET：myapi.Api(url,myapi.MYAPI_METHOD_GET,map_data)
- POST-json：myapi.Api(url,myapi.MYAPI_METHOD_JSON,map_data)
- POST-form：myapi.Api(url,myapi.MYAPI_METHOD_POST_FORM,map_data)
- 不记录请求和响应日志：myapi.ApiNoAllLog
- 不记录响应日志：myapi.ApiNoResponsLog
### Mysql
```go
var obj entity.Ademo
orm := core.Db
orm.Where("name=?", name).Find(&obj)
```
### Redis
- 封装后的写法（仅get set del，请求方法不全，可以自行完善）
```go
myredis := MyRedis{}
defer myredis.Close()
myredis.Set(key, val, ttl)
val,err:=myredis.Get(key)
if myredis.IsError(err) {
	return err happend...
}
```
- 封装conn写法-原生请求
```go
myredis := MyRedis{}
defer myredis.Close()
conn := myredis.GetConn()
conn.Do("Set", "abc", 100, "EX", 100)
```
- 原生写法
```go
conn := core.Redis_pool.Get()
defer conn.Close()
_, err := conn.Do("Set", "abc", 100, "EX", 100)

res, err := redis.Int(conn.Do("Get", "abc"))
```
## 其他
- 默认控制器为Home，默认action为Index，见配置文件router
- 控制器中的CallBefore为钩子方法，在运行acton之前运行，return false可阻断运行
- 接口返回推荐：core.ResultFail(ctx,msg),core.ResultSucc(ctx,data)
- 项目常量位置：consts/user_consts.go
- 用户库位置：lib/userlib/
- 判断运行环境：core.IsEnvDev() core.IsEnvTest() core.IsEnvPro()
- 数组方法：lib/helper/arrmap/arr.go InArrayStr() InArrayInt() IsArray()...
- map方法：lib/helper/arrmap/map.go InMapStrKey() InMapStrVal() MergeMap()...
- 字符串方法：lib/helper/str/str.go FirstToUpper() Md5() ToString() Json_decode() TimeFormat()...
- 正则验证：lib/helper/str/verify/verify.go IsDate() IsTime() IsNumber()...

## 开始编码
比如： ./user/login
代码初始化
- 复制 ./controller/AdemoController.go =》UserController.go 
- 复制 ./model/AdemoModel.go =》UserModel.go
- 复制 ./service/AdemoService.go =》 UserService.go
- 复制 ./entity/Ademo.go =》 User.go
- 修改上述文件中的type名称，如type AdemoController struct等
- 添加控制器映射：./mygin/router.go->InitCtrlByName方法中增加
    case "user":
		v = reflect.ValueOf(new(controller.UserController))
- 开始编码：UserController.go -> login()
		
## 启动运行
支持start/stop/restart
- 开发环境：go run main.go start/stop/restart
- 生产环境：xxx start/stop/restart &

## 联系我
- 409250643@qq.com

