package global

import (
	"github.com/xormplus/xorm"
	"log"
	"sync"
)

var (
	Engin *xormEngin
)

type XormInterface interface {
	SetXormEngin(k string, e *xorm.Engine)
	GetXormEngin(k string) (e *xorm.Engine)
}

type XormEngin struct {
	*xormEngin
}

type xormEngin struct {
	mux   sync.Mutex
	items map[string]*xorm.Engine
}

func init() {
	Engin = NewXormPool()
	for _, db := range Config.Db {
		engine, err := db.GetEngin()
		if err != nil {
			log.Fatalln(err)
		}
		Engin.SetXormEngin(db.Name, engine)
	}
}

func NewXormPool() *xormEngin {
	x := &xormEngin{}
	x.items = make(map[string]*xorm.Engine)
	return x
}

func (x *xormEngin) SetXormEngin(k string, e *xorm.Engine) {
	x.mux.Lock()
	x.items[k] = e
	x.mux.Unlock()
}

func (x *xormEngin) GetXormEngin(k string) (e *xorm.Engine, found bool) {
	x.mux.Lock()
	defer x.mux.Unlock()
	if e, found = x.items[k]; !found {
		return nil, false
	}
	return
}
