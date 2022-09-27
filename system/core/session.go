package core

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
)

var Session_store *redistore.RediStore
var session_name string

const ROLLING_TIME_SESS_KEY = "__rolling_time__"

var ROLLING_INTERVAL_SECOND = time.Second.Nanoseconds() * 60 //刷新间隔时间，默认60(秒)，刚刷过60秒内不再刷

/*
*
用法1：

	session, err := core.InitSession(ctx.Request, ctx.Writer,entity.user{})
	if err != nil {
		//err
	}

	//保存值
	err = session.Save("foo", time.Now().Local().String())

	//取值
	foo := session.Get("foo")

	//清空session-退出
	err = session.Del()

用法2：
默认一个固定的key,const SESS_USER_KEY = "gin_user"

	//保存
	err:=core.SaveUserSess(ctx,val,entity.user{})
	//取值
	val,err:=err:=GetUserSess(ctx,,entity.user{})
	//清空退出
	err:=DelUserSess(ctx)
*/
func init() {
	session_name = GetConfigString("session.name")
	// Note: Don't store your key in your source code. Pass it via an
	// environmental variable, or flag (or both), and don't accidentally commit it
	// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
	// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
	store_key := GetConfigString("session.store_key")
	max_age := GetConfigInt("session.max_age")
	if session_name == "" || store_key == "" || max_age < 1 {
		LogInfo("session.name,session.store_key,session.max_age未配置，session未初始化")
		return
	}

	//store = sessions.NewCookieStore([]byte(store_key)) //[]byte(os.Getenv("SESSION_KEY"))
	var err error
	Session_store, err = redistore.NewRediStoreWithPool(Redis_pool, []byte(store_key))
	if err != nil {
		panic(err)
	}

	//defer store.Close()
	Session_store.SetMaxAge(max_age)

	LogInfo("init session")
}

type Session struct {
	Sess      *sessions.Session
	req       *http.Request
	rpw       http.ResponseWriter
	IsRolling bool //是否每次请求都刷新过期时间，默认true
	isFresh   bool //是否刷新过有效期，一次请求只刷新一次
}

// 初试session-(配置文件的session.name)
func InitSession(req *http.Request, rpw http.ResponseWriter, reg_vals ...interface{}) (*Session, error) {
	for _, reg_val := range reg_vals {
		gob.Register(reg_val)
	}

	sess, err := InitSessionBySessname(req, rpw, session_name)
	if err != nil {
		LogError("初试化session失败，" + err.Error())
	}
	return sess, err
}

// 初试session-指定session_name
func InitSessionBySessname(req *http.Request, rpw http.ResponseWriter, sess_name string) (*Session, error) {
	sess := Session{}
	sess.IsRolling = true
	sess.isFresh = false
	err := sess.getSessionBySessname(req, rpw, sess_name)
	return &sess, err
}

// 获得指定session-指定session_name
func (self *Session) getSessionBySessname(req *http.Request, rpw http.ResponseWriter, sess_name string) error {
	if session_name == "" {
		return errors.New("session未初始化")
	}
	sess, err := Session_store.Get(req, sess_name)
	if err != nil {
		return err
	}
	sess.Options.HttpOnly = true
	sess.Options.Secure = true
	self.Sess = sess
	self.req = req
	self.rpw = rpw
	return nil
}

// 清除默认session
func (self *Session) Del() error {
	self.Sess.Options.MaxAge = -1
	err := self.Sess.Save(self.req, self.rpw)
	if err != nil {
		LogError("del session失败，" + err.Error())
	}
	return err
}

// 获取session内容
func (self *Session) Get(session_key string) interface{} {
	if self.IsRolling && !self.isFresh { //获取的时候顺便刷新下过期时间
		last_roll_time, has := self.Sess.Values[ROLLING_TIME_SESS_KEY]
		if !has || (last_roll_time.(int64) < (time.Now().UnixNano() - ROLLING_INTERVAL_SECOND)) { //不存在，或超过了60秒，就刷新
			self.isFresh = true
			self.Save(ROLLING_TIME_SESS_KEY, time.Now().UnixNano()) //每次请求都刷新过期时间。。
			//self.SaveNoWithVal()
		}
	}

	return self.Sess.Values[session_key]
}

// 保存session
func (self *Session) Save(session_key string, session_val interface{}) error {
	self.Sess.Values[session_key] = session_val
	err := self.Sess.Save(self.req, self.rpw)
	if err != nil {
		LogError("save session失败，" + err.Error())
	}
	return err
}
