package main

import (
	"bytes"
	"./config"
)

// TrainingMessageBuilder helps you to build messages
type TrainingMessageBuilder struct {
	config *config.Config
}

func newTrainingMessageBuilder(config *config.Config) *TrainingMessageBuilder {
	messageBuilder := new(TrainingMessageBuilder)
	messageBuilder.config = config

	return messageBuilder
}

func (messageBuilder TrainingMessageBuilder) create(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString(row[messageBuilder.config.NameColumn].(string))
	buffer.WriteString(" am *")
	buffer.WriteString(row[messageBuilder.config.DayOfWeekColumn].(string))
	buffer.WriteString(", ")
	buffer.WriteString(row[messageBuilder.config.DateColumn].(string))
	buffer.WriteString("* Uhr! ")
	buffer.WriteString(row[messageBuilder.config.DescriptionColumn].(string))

	return buffer
}
