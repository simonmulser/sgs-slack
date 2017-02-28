package main

import (
	"bytes"
	"time"

	"github.com/golang/glog"
)

type IMessageBuilder interface {
	create(row []interface{}) bytes.Buffer
}

// MessageBuilder helps you to build messages
type MessageBuilder struct {
	config *Config
	ISlackService
}

func newMessageBuilder(config *Config, slackService ISlackService) *MessageBuilder {
	messageBuilder := new(MessageBuilder)
	messageBuilder.ISlackService = slackService
	messageBuilder.config = config

	return messageBuilder
}

type trainingParameters struct {
	TotalGoing                  string
	GoingSGS07                  string
	GoingSGS16                  string
	ResponsibleTrainingUtensils string
}

func (messageBuilder MessageBuilder) createTrainingPost(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(row[messageBuilder.config.NameColumn].(string))
	buffer.WriteString(" am *")
	buffer.WriteString(row[messageBuilder.config.DateColumn].(string))
	buffer.WriteString("* Uhr! ")
	buffer.WriteString(row[messageBuilder.config.DescriptionColumn].(string))
	return buffer
}

func (messageBuilder MessageBuilder) create(row []interface{}) bytes.Buffer {
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
