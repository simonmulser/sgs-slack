package main

import (
	"time"

	"github.com/golang/glog"
)

func timeNow() time.Time {
	t := time.Now()
	utc, error := time.LoadLocation("Europe/Vienna")
	if error != nil {
		glog.Info("error: ", error.Error())
	}
	return t.In(utc)
}
