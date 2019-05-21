package model

import (
    "time"
)

type ProductFile struct {
    Id            int       `xorm:"not null pk autoincr unique INT(11)"`
    ProductTypeId int       `xorm:"not null INT(11)"`
    Name          string    `xorm:"not null TEXT"`
    Url           string    `xorm:"not null TEXT"`
    CompressUrl1  string    `xorm:"not null TEXT"`
    CompressUrl2  string    `xorm:"not null TEXT"`
    CreateTime    time.Time `xorm:"not null TIMESTAMP created"`
    UpdateTime    time.Time `xorm:"not null TIMESTAMP updated"`
    IsDelete      int       `xorm:"not null default 0 INT(11)"`
}
