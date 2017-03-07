package main

import (
	"errors"
	"testing"

	"github.com/nlopes/slack"
	"github.com/simonmulser/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	config := testutils.CreateConfig()
	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	itemRef := slack.ItemRef{Channel: "channelID", Timestamp: "Timestamp"}
	mockSlackService.On("GetReactions", itemRef, slack.GetReactionsParameters{false}).Return([]slack.ItemReaction{}, nil)
	mockSlackService.On("PostMessage", mock.MatchedBy(func(s []string) bool { return true })).Return("channelID", "timestamp", nil)
	mockSpreadsheetService := new(MockSpreadsheetService)
	mockSpreadsheetService.On("WriteCell", "", 0, config.TrainingUtensilsColumn, "POSTED").Return()

	mockTrainingParamsService := new(MockTrainingParamsService)
	mockTrainingParamsService.On("create", mock.MatchedBy(func(r []slack.ItemReaction) bool { return true })).Return(trainingParameters{"", "", "", ""})

	trainingService := newTrainingService(config, mockSlackService, mockSpreadsheetService, mockTrainingParamsService)
	row := testutils.CreateRow([]string{"POSTED", "05.05.1991 20:20", "channelID", "Timestamp", "05.05.1991 20:20", "NOT_POSTED"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.Nil(t, error)
}

func TestExecuteErrorFromSlack(t *testing.T) {

	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	itemRef := slack.ItemRef{Channel: "channelID", Timestamp: "Timestamp"}
	mockSlackService.On("GetReactions", itemRef, slack.GetReactionsParameters{false}).Return([]slack.ItemReaction{}, errors.New("errorFromMock"))
	mockSpreadsheetService := new(MockSpreadsheetService)
	mockTrainingParamsService := new(MockTrainingParamsService)

	trainingService := newTrainingService(testutils.CreateConfig(), mockSlackService, mockSpreadsheetService, mockTrainingParamsService)
	row := testutils.CreateRow([]string{"POSTED", "", "channelID", "Timestamp", "05.05.1991 20:20", "NOT_POSTED"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
}

func TestExecuteWrongDate(t *testing.T) {

	topicConfig, _ := createTopicConfig()
	mockSlackService := new(MockSlackService)
	mockSpreadsheetService := new(MockSpreadsheetService)
	mockTrainingParamsService := new(MockTrainingParamsService)

	trainingService := newTrainingService(testutils.CreateConfig(), mockSlackService, mockSpreadsheetService, mockTrainingParamsService)
	row := testutils.CreateRow([]string{"", "", "", "", "05.05.1991"})

	error := trainingService.execute(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}
