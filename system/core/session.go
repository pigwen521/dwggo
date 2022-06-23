package core

import (
	"errors"
	"net/http"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
)

var Session_store *redistore.RediStore
var session_name string

/**
用法：
	session, err := core.InitSession(ctx.Request, ctx.Writer)
	if err != nil {
		//err
	}

	//保存值
	err = session.Save("foo", time.Now().Local().String())

	//取值
	foo := session.Get("foo")

	//清空session-退出
	err = session.Del()
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
	Sess *sessions.Session
	req  *http.Request
	rpw  http.ResponseWriter
}

//初试session-(配置文件的session.name)
func InitSession(req *http.Request, rpw http.ResponseWriter) (*Session, error) {
	return InitSessionBySessname(req, rpw, session_name)
}

//初试session-指定session_name
func InitSessionBySessname(req *http.Request, rpw http.ResponseWriter, sess_name string) (*Session, error) {
	sess := Session{}
	err := sess.getSessionBySessname(req, rpw, sess_name)
	return &sess, err
}

//获得指定session-指定session_name
func (self *Session) getSessionBySessname(req *http.Request, rpw http.ResponseWriter, sess_name string) error {
	if session_name == "" {
		return errors.New("session未初始化")
	}
	sess, err := Session_store.Get(req, sess_name)
	if err != nil {
		return err
	}
	self.Sess = sess
	self.req = req
	self.rpw = rpw
	return nil
}

//清除默认session
func (self *Session) Del() error {
	self.Sess.Options.MaxAge = -1
	return self.Sess.Save(self.req, self.rpw)
}

//获取session内容
func (self *Session) Get(session_key string) interface{} {
	return self.Sess.Values[session_key]
}

//保存session
func (self *Session) Save(session_key string, session_val interface{}) error {
	self.Sess.Values[session_key] = session_val
	return self.Sess.Save(self.req, self.rpw)
}
