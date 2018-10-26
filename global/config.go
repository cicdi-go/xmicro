package global

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/xormplus/xorm"
	"log"
	"os"
	"path/filepath"
)

var (
	Config Conf
)

type Conf struct {
	App            string     `yaml:"app"`
	Secret         string     `yaml:"secret"`
	SsoTokenExpire int64      `yaml:"sso-token-expire"`
	TablePrefix    string     `yaml:"table-prefix"`
	Db             []DateBase `yaml:"db"`
}

type DateBase struct {
	Name         string `yaml:"name"`
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	Log          string `yaml:"log"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	ShowSql      bool   `yaml:"show-sql"`
}

func init() {
	if err := configor.Load(&Config, "config/app.yml", "config/app-local.yml"); err != nil {
		log.Fatal(err)
	}
}

func (db *DateBase) GetEngin() (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine(db.Driver, db.Dsn)
	if err != nil {
		return engine, err
	}
	if db.Log != "" {
		dir := filepath.Dir(db.Log)
		exist, err := PathExists(dir)
		if err != nil {
			fmt.Printf("get dir error![%v]\n", err)
			return engine, err
		}
		if !exist {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
		var f *os.File
		f, err = os.Create(db.Log)
		if err != nil {
			return engine, err
		}
		engine.SetLogger(xorm.NewSimpleLogger(f))
		engine.ShowSQL(true)
	}
	if db.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(db.MaxIdleConns)
	}
	if db.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(db.MaxOpenConns)
	}
	return
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
