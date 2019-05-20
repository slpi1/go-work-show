package service

import (
	"time"
    "sync"

	"lib"
)

type Timer struct {
	TimeLine map[string] time.Time
}

var timer *Timer
var once sync.Once

func NewTimer() *Timer {
	once.Do(func() {
		timer = new(Timer)
		timer.TimeLine = make(map[string] time.Time)
    })
	return timer
}

func (timer *Timer) Start() {
	timer.TimePoint("start")
}

func (timer *Timer) End() {
	timer.TimePoint("end")
}

func (timer *Timer) TimePoint(name string) time.Time {
	t := time.Now()
	timer.TimeLine[name] = t
	lib.Logger().Println(name, t.String())
	return t
}

func (timer *Timer) Diff(p1 string, p2 string) string {
	time1, ok := timer.TimeLine[p1] 
	if !ok {
		return ""
	}
	time2, ok := timer.TimeLine[p2] 
	if !ok {
		return ""
	}
	length := time2.Sub(time1).String()
	lib.Logger().Println(p1, "->" ,p2, ":", length)
	return length
}