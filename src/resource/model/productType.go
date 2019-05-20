package model

import (
	"time"
)

type ProductType struct {
	Id             		int       `xorm:"not null pk autoincr unique INT(11)"`
	SupplierId     		int       `xorm:"not null INT(11)"`
	Name           		string    `xorm:"not null TEXT"`
	Url            		string    `xorm:"not null TEXT"`
	Fraction       		float32   `xorm:"not null default 0 FLOAT"`
	CooperationNum 		int       `xorm:"not null default 0 INT(11)"`
	ParentId       		int       `xorm:"not null default -1 INT(11)"`
	CoverPath      		string    `xorm:"not null TEXT"`
    CoverCompressPath1 	string    `xorm:"not null TEXT"`
    CoverCompressPath2 	string    `xorm:"not null TEXT"`	
	CreateTime     		time.Time `xorm:"not null TIMESTAMP created"`
	UpdateTime     		time.Time `xorm:"not null TIMESTAMP updated"`
	IsDelete       		int       `xorm:"not null default 0 INT(11)"`
}
