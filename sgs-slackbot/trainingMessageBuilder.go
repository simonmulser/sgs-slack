package main

import (
	"bytes"

	"github.com/nlopes/slack"
)

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

func (messageBuilder TrainingBuilder) createTrainingPost(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer

	return buffer
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

func (messageBuilder TrainingBuilder) createTrainingMgmtPost(row []interface{}, params trainingParameters) bytes.Buffer {
	var buffer bytes.Buffer
	return buffer
}

func (messageBuilder TrainingBuilder) createTrainingParams(reactions []slack.ItemReaction) trainingParameters {
	var params trainingParameters

	return params
}
