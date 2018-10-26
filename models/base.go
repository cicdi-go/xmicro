package models

import (
	"errors"
	"github.com/cicdi-go/xmicro/global"
	"github.com/xormplus/xorm"
)

type ActiveRecod interface {
	GetDb() (e *xorm.Engine, err error)
	TableName() string
}

type Base struct {
}

func (u *Base) TableName() string {
	return global.Config.TablePrefix + "user"
}

func (u *Base) GetDb() (e *xorm.Engine, err error) {
	var found bool
	if e, found = global.Engin.GetXormEngin("default"); !found {
		err = errors.New("Database default is not found!")
	}
	return
}
