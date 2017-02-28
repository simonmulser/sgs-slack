package main

import (
	"bytes"
	"time"

	"github.com/golang/glog"
)

// GameMessageBuilder helps you to build messages
type GameMessageBuilder struct {
	config *Config
}

func newGameMessageBuilder(config *Config) *GameMessageBuilder {
	messageBuilder := new(GameMessageBuilder)
	messageBuilder.config = config

	return messageBuilder
}

func (messageBuilder GameMessageBuilder) create(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer

	meetingTime, error := time.Parse("02.01.2006 15:04", row[messageBuilder.config.DateColumn].(string))
	if error != nil {
		glog.Fatalf("Unable to parse date. %v", error)
	}
	meetingTime = meetingTime.Add(-1 * 60 * time.Minute)

	buffer.WriteString("*")
	buffer.WriteString(row[messageBuilder.config.HomeColumn].(string))
	buffer.WriteString(" : ")
	buffer.WriteString(row[messageBuilder.config.AwayColumn].(string))
	buffer.WriteString("* am *")
	buffer.WriteString(row[messageBuilder.config.DateColumn].(string))
	buffer.WriteString("* Uhr auf ")
	buffer.WriteString(row[messageBuilder.config.SurfaceColumn].(string))
	buffer.WriteString("! Treffpunkt: ")
	buffer.WriteString(meetingTime.Format("15:04"))
	buffer.WriteString(", ")
	buffer.WriteString(row[messageBuilder.config.LocationColumn].(string))
	buffer.WriteString(" - ")
	buffer.WriteString(row[messageBuilder.config.LocationMapsColumn].(string))
	buffer.WriteString(".")

	return buffer
}
