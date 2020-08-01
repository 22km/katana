package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// Config ...
// "conf_keeper_twik01_rw:yOF0bPHxsAt1eRz@tcp(10.88.128.43:4000)/conf_keeper?charset=utf8&autocommit=true"
// dataSourceName = "root:123456@tcp(10.96.80.244:3306)/conf_keeper?charset=utf8&autocommit=true"
type Config struct {
	User         string `yaml:"user"`
	Pwd          string `yaml:"pwd"`
	IP           string `yaml:"ip"`
	Port         string `yaml:"port"`
	DB           string `yaml:"db"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
}

var engine *xorm.Engine

// Init ...
func Init(c *Config) error {
	src := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&autocommit=true",
		c.User,
		c.Pwd,
		c.IP,
		c.Port,
		c.DB)

	// fmt.Println(src)

	e, err := xorm.NewEngine("mysql", src)
	if err != nil {
		return err
	}
	e.SetMaxIdleConns(c.MaxIdleConns)
	e.ShowSQL(true)

	engine = e
	return nil
}

// Engine ...
func Engine() *xorm.Engine {
	return engine
}
