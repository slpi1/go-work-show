package model

import (
	"time"
)

type Reservation struct {
	Id            int       `xorm:"not null pk autoincr unique INT(11)"`
	SupplierId    string    `xorm:"not null TEXT"`
	ProductTypeId string    `xorm:"not null TEXT"`
	ProjectName   string    `xorm:"not null VARCHAR(45)"`
	DemandDate    time.Time `xorm:"TIMESTAMP"`
	ClearDemand   int       `xorm:"not null INT(11)"`
	Instruction   string    `xorm:"not null TEXT"`
	ReceiverList  string    `xorm:"not null TEXT"`
	Founder       int       `xorm:"not null INT(11)"`
	CreateTime    time.Time `xorm:"not null TIMESTAMP created"`
	UpdateTime    time.Time `xorm:"not null TIMESTAMP updated"`
	IsDelete      string    `xorm:"not null default '0' VARCHAR(45)"`
}
