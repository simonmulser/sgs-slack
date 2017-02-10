package main

import (
	"time"

	"github.com/golang/glog"
)

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func toAlphabetChar(i int) string {
	return arr[i]
}

func timeNow() time.Time {
	t := time.Now()
	utc, error := time.LoadLocation("Europe/Vienna")
	if error != nil {
		glog.Info("error: ", error.Error())
	}
	return t.In(utc)
}
