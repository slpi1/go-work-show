package lib

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func Connection() *xorm.Engine{
	if engine != nil {
		return engine
	}

    config := NewConfig()

	var dsn = config.Db.Username + ":"+config.Db.Password+"@"+config.Db.Url+"?charset=utf8"
	var err error
    engine, err = xorm.NewEngine("mysql", dsn)
	engine.SetMaxOpenConns(50)
    if err != nil {
    	return nil
    }
    return engine
}