package model

import (
	"time"
)

type Recommended struct {
	Id           int       `xorm:"not null pk autoincr unique INT(11)"`
	SupplierId   int       `xorm:"not null INT(11)"`
	SupplierType int       `xorm:"not null INT(11)"`
	Reason       string    `xorm:"TEXT"`
	Founder      int       `xorm:"not null INT(11)"`
	Updater      int       `xorm:"not null INT(11)"`
	CreateTime   time.Time `xorm:"not null TIMESTAMP created"`
	UpdateTime   time.Time `xorm:"not null TIMESTAMP updated"`
	IsDelete     int       `xorm:"not null default 0 INT(11)"`
}
