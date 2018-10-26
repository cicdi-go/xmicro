package models

import (
	"fmt"
	"github.com/cicdi-go/xmicro/global"
	"github.com/pkg/errors"
	"log"
)

type User struct {
	*Base              `xorm:"-"`
	Id                 int64  `json:"id"`
	Username           string `xorm:"varchar(100) notnull unique index default ''" json:"username"`
	Status             int    `xorm:"SMALLINT default 1" json:"status"`
	AuthKey            string `xorm:"varchar(32) default ''" json:"-"`
	PasswordHash       string `xorm:"varchar(255) default ''" json:"-"`
	PasswordResetToken string `xorm:"varchar(255) default ''" json:"-"`
	password           string `xorm:"-"`
	CreatedAt          int    `xorm:"created" json:"created_at"`
	UpdatedAt          int    `xorm:"updated" json:"updated_at"`
}

func (u *User) TableName() string {
	return global.Config.TablePrefix + "user"
}

func init() {
	u := new(User)
	if e, err := u.GetDb(); err != nil {
		log.Println(err)
	} else {
		err := e.Sync2(u)
		if err != nil {
			log.Println(err)
		}
	}
}

func (u *User) Exist() (exist bool, err error) {
	engine, err := u.GetDb()
	if err != nil {
		return
	}
	return engine.Exist(u)
}

func (u *User) Insert() (err error) {
	engine, err := u.GetDb()
	if err != nil {
		return
	}
	id, err := engine.Insert(u)
	if err != nil {
		return err
	}
	u.Id = id
	return
}

func (u *User) SetPassword(value string) {
	u.password = value
	u.generateAuthKey()
	u.PasswordHash, _ = global.SetPassword(u.password, u.AuthKey)
}

func (u *User) GetPasswordHash(p string) string {
	passwordHash, err := global.SetPassword(p, u.AuthKey)
	if err != nil {
		return ""
	}
	return passwordHash
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) generateAuthKey() {
	u.AuthKey = global.GenerateRandomKey()
}

func (u *User) Verify(p string) bool {
	engine, err := u.GetDb()
	if err != nil {
		return false
	}
	engine.Where("username = ?", u.Username).Get(u)
	return u.GetPasswordHash(p) == u.PasswordHash
}

func FindByUsername(username string) (u User, err error) {
	u.Username = username
	engine, err := u.GetDb()
	if err != nil {
		return
	}
	if has, err := engine.Get(&u); err != nil {
		return u, err
	} else if !has {
		return u, errors.New(fmt.Sprintf("%s is not exist!", u.Username))
	}
	return
}
