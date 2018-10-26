package services

import (
	"context"
	"github.com/cicdi-go/xmicro/global"
	"github.com/cicdi-go/xmicro/models"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/robbert229/jwt"
	"time"
)

type UserRegister struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type User int

type Reply struct {
	Status int `json:"status"`
}

type Jwt struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

//用户注册
func (u *User) Registry(ctx context.Context, arge *UserRegister, reply *Reply) (err error) {
	validate := validator.New()
	if err = validate.Struct(arge); err != nil {
		return
	}

	user := &models.User{
		Username: arge.Username,
	}
	if exist, err := user.Exist(); err != nil {
		return err
	} else if exist {
		return errors.New("Username already exists")
	}
	user.Status = global.STATUS_ENABLE
	user.SetPassword(arge.Password)
	if err = user.Insert(); err != nil {
		return
	}
	reply.Status = user.Status
	return
}

//获取token
func (u *User) Authorization(ctx context.Context, arge *UserRegister, reply *Jwt) (err error) {
	user, err := models.FindByUsername(arge.Username)
	if err != nil {
		return err
	}
	if user.Verify(arge.Password) {
		algorithm := jwt.HmacSha256(global.Config.Secret)

		claims := jwt.NewClaim()
		claims.Set("Username", arge.Username)
		claims.SetTime("exp", time.Now().Add(time.Duration(global.Config.SsoTokenExpire)*time.Minute))

		if token, err := algorithm.Encode(claims); err != nil {
			return err
		} else {
			reply.Token = token
			reply.Expire = time.Now().Add(time.Duration(global.Config.SsoTokenExpire) * time.Minute)
		}
	} else {
		return errors.New("Verification failed")
	}
	return
}
