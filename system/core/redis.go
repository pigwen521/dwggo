package core

import (
	"errors"
	"time"

	//"github.com/garyburd/redigo/redis"
	"github.com/gomodule/redigo/redis"
)

var Redis_pool *redis.Pool

/**
1,
myredis := core.MyRedis{}
defer myredis.Close()
myredis.Set(key, val, ttl)
val,err:=myredis.Get(key)
if myredis.IsError(err) {
	return err happend...
}
2,
myredis := core.MyRedis{}
defer mr.Close()
conn := myredis.GetConn()
conn.Do("Set", "abc", 100, "EX", 100)
3,
conn := core.Redis_pool.Get()
defer conn.Close()
_, err := conn.Do("Set", "abc", 100, "EX", 100)

res, err := redis.Int(conn.Do("Get", "abc"))
*/
func init() {
	host := GetConfigString("redis.host")
	port := GetConfigString("redis.port")
	password := redis.DialPassword(GetConfigString("redis.password"))
	Redis_pool = &redis.Pool{ //实例化一个连接池
		MaxIdle:     GetConfigInt("redis.MaxIdle"),                          //最大空闲连接
		MaxActive:   GetConfigInt("redis.MaxActive"),                        //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: time.Second * (GetConfigDuration("redis.IdleTimeout")), //连接关闭时间 300秒 （300秒不使用自动关闭）
		//MaxConnLifetime: time.Second * ,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host+":"+port, password)
		},
	}

	LogInfo("init redis pool")
}
func RedisPoolClose() {
	Redis_pool.Close()
}

type MyRedis struct {
	Conn redis.Conn
}

func (self *MyRedis) GetConn() redis.Conn {
	if self.Conn == nil {
		self.Conn = Redis_pool.Get()
	}
	return self.Conn
}

//ttl_second 小于等于0，为永久有效
func (self *MyRedis) Set(key string, val interface{}, ttl_second int) (interface{}, error) {
	self.GetConn()
	if ttl_second <= 0 {
		return self.Conn.Do("Set", key, val)
	} else {
		return self.Conn.Do("Set", key, val, "EX", ttl_second)
	}
}
func (self *MyRedis) Get(key string) (interface{}, error) {
	self.GetConn()
	return self.Conn.Do("Get", key)
}
func (self *MyRedis) Del(key string) (interface{}, error) {
	self.GetConn()
	return self.Conn.Do("Del", key)
}

//判断Get后的是否有error，要排查为空的数据
func (self *MyRedis) IsError(err error) bool {
	return err != nil && !errors.Is(err, redis.ErrNil)
}

func (self *MyRedis) Close() {
	self.Conn.Close()
}
