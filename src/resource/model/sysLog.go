package model

import (
    "time"
)

type SysLog struct {
    Id              int64     `xorm:"pk autoincr unique BIGINT(20)"`
    LogTime         time.Time `xorm:"not null TIMESTAMP created"`
    LogLevel        string    `xorm:"not null VARCHAR(30)"`
    LogPosition     string    `xorm:"VARCHAR(100)"`
    ClientIpAddress int       `xorm:"INT(10)"`
    ServerIpAddress int       `xorm:"INT(10)"`
    AppName         string    `xorm:"VARCHAR(50)"`
    Context         string    `xorm:"not null TEXT"`
    Extend1         string    `xorm:"VARCHAR(64)"`
    Extend2         string    `xorm:"VARCHAR(128)"`
}
