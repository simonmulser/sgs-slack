package main

import "bytes"

// TrainingBuilder helps you to build messages
type TrainingBuilder struct {
	config *Config
	ISlackService
}

func newTrainingBuilder(config *Config, slackService ISlackService) *TrainingBuilder {
	messageBuilder := new(TrainingBuilder)
	messageBuilder.ISlackService = slackService
	messageBuilder.config = config

	return messageBuilder
}

func (messageBuilder TrainingBuilder) create(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString(row[messageBuilder.config.NameColumn].(string))
	buffer.WriteString(" am *")
	buffer.WriteString(row[messageBuilder.config.DateColumn].(string))
	buffer.WriteString("* Uhr! ")
	buffer.WriteString(row[messageBuilder.config.DescriptionColumn].(string))

	return buffer
}
