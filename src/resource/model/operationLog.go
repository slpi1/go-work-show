package model

import (
    "time"
)

type OperationLog struct {
    Id      int       `xorm:"not null pk autoincr unique INT(11)"`
    LogTime time.Time `xorm:"not null TIMESTAMP created"`
    LogType string    `xorm:"not null VARCHAR(45)"`
    Content string    `xorm:"not null TEXT"`
    Founder int       `xorm:"not null INT(11)"`
}
