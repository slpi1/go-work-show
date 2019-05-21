package model

import (
    "time"
)

type Supplier struct {
    Id                 int       `xorm:"not null pk autoincr unique INT(11)"`
    Name               string    `xorm:"not null unique TEXT"`
    Code               string    `xorm:"not null unique VARCHAR(10)"`
    Url                string    `xorm:"not null TEXT"`
    Type               int       `xorm:"not null INT(11)"`
    CoverPath          string    `xorm:"not null TEXT"`
    CoverCompressPath1 string    `xorm:"not null TEXT"`
    CoverCompressPath2 string    `xorm:"not null TEXT"`
    CreateTime         time.Time `xorm:"not null TIMESTAMP created"`
    UpdateTime         time.Time `xorm:"not null TIMESTAMP updated"`
    IsDelete           int       `xorm:"not null default 0 INT(11)"`
}
