package main

import (
	"errors"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	config := createConfig()
	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	itemRef := slack.ItemRef{Channel: "channelID", Timestamp: "Timestamp"}
	mockSlackService.On("getReactions", itemRef, slack.GetReactionsParameters{false}).Return([]slack.ItemReaction{}, nil)
	mockSlackService.On("postMessage", mock.MatchedBy(func(s []string) bool { return true })).Return("channelID", "timestamp", nil)
	mockSpreadsheetService := new(MockSpreadsheetService)
	mockSpreadsheetService.On("writeCell", "", 0, config.TrainingUtensilsColumn, "POSTED").Return()

	trainingService := newTrainingService(config, mockSlackService, mockSpreadsheetService)
	row := createRow([]string{"POSTED", "05.05.1991 20:20", "channelID", "Timestamp", "05.05.1991 20:20", "NOT_POSTED"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.Nil(t, error)
}

func TestExecuteErrorFromSlack(t *testing.T) {

	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	itemRef := slack.ItemRef{Channel: "channelID", Timestamp: "Timestamp"}
	mockSlackService.On("getReactions", itemRef, slack.GetReactionsParameters{false}).Return([]slack.ItemReaction{}, errors.New("errorFromMock"))
	mockSpreadsheetService := new(MockSpreadsheetService)

	trainingService := newTrainingService(createConfig(), mockSlackService, mockSpreadsheetService)
	row := createRow([]string{"POSTED", "", "channelID", "Timestamp", "05.05.1991 20:20", "NOT_POSTED"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
}

func TestExecuteWrongDate(t *testing.T) {

	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	mockSpreadsheetService := new(MockSpreadsheetService)

	trainingService := newTrainingService(createConfig(), mockSlackService, mockSpreadsheetService)
	row := createRow([]string{"", "", "", "", "05.05.1991"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}
