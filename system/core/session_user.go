package core

import (
	"encoding/gob"

	"github.com/gin-gonic/gin"
)

const SESS_USER_KEY = "gin_user"

//保存-key默认gin_user
func SaveUserSess(ctx *gin.Context, val interface{}, reg_vals ...interface{}) error {
	for _, reg_val := range reg_vals {
		gob.Register(reg_val)
	}
	session, err := InitSession(ctx.Request, ctx.Writer)
	if err != nil {
		return err
	}
	//保存值
	return session.Save(SESS_USER_KEY, val)
}

//取值-key默认gin_user
func GetUserSess(ctx *gin.Context, reg_vals ...interface{}) (interface{}, error) {
	for _, reg_val := range reg_vals {
		gob.Register(reg_val)
	}
	session, err := InitSession(ctx.Request, ctx.Writer)
	if err != nil {
		return nil, err
	}

	return session.Get(SESS_USER_KEY), nil
}

//清空session-退出-会清空全部，不仅仅gin_user
func DelUserSess(ctx *gin.Context) error {
	session, err := InitSession(ctx.Request, ctx.Writer)
	if err != nil {
		return err
	}
	return session.Del()
}
