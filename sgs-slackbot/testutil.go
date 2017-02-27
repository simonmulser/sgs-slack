package main

import (
	"bytes"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

type MockSpreadsheetService struct {
	mock.Mock
}

func (mockSpreadsheetService MockSpreadsheetService) writeCell(sheet string, row int, column int, text string) {
	mockSpreadsheetService.Called(sheet, row, column, text)
	return
}

type MockSlackService struct {
	mock.Mock
}

func (mockSlackService MockSlackService) postMessage(channel string, message string) (string, string, error) {
	args := mockSlackService.Called(channel, message)
	return args.String(0), args.String(1), args.Error(2)
}

func (mockSlackService MockSlackService) updateMessage(channel, timestamp, text string) (string, string, string, error) {
	args := mockSlackService.Called(channel, timestamp, text)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

type MockMessageBuilder struct {
	mock.Mock
}

func (messageBuilder MockMessageBuilder) createTrainingPost(row []interface{}) bytes.Buffer {
	var buffer bytes.Buffer
	return buffer
}

func (messageBuilder MockMessageBuilder) createGamePost(row []interface{}) bytes.Buffer {
	args := messageBuilder.Called(row)
	return args.Get(0).(bytes.Buffer)
}

func (messageBuilder MockMessageBuilder) createTrainingMgmtPost(row []interface{}, params trainingParameters) bytes.Buffer {
	args := messageBuilder.Called(row, params)
	return args.Get(0).(bytes.Buffer)
}

func (messageBuilder MockMessageBuilder) createTrainingParams(reactions []slack.ItemReaction) trainingParameters {
	args := messageBuilder.Called(reactions)
	return args.Get(0).(trainingParameters)
}
