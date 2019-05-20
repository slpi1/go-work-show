package model

import (
	"time"
)

type SupplierType struct {
	Id         int       `xorm:"not null pk autoincr unique INT(11)"`
	Name       string    `xorm:"not null unique VARCHAR(45)"`
	Url        string    `xorm:"not null TEXT"`
	CreateTime time.Time `xorm:"not null TIMESTAMP created"`
	UpdateTime time.Time `xorm:"not null TIMESTAMP updated"`
	IsDelete   int       `xorm:"not null default 0 INT(11)"`
}
