package main

import (
	"bytes"

	sheets "google.golang.org/api/sheets/v4"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

type MockTrainingParamsService struct {
	mock.Mock
}

func (mockTrainingParamsService MockTrainingParamsService) create(reactions []slack.ItemReaction) trainingParameters {
	args := mockTrainingParamsService.Called(reactions)
	return args.Get(0).(trainingParameters)
}

type MockSpreadsheetService struct {
	mock.Mock
}

func (mockSpreadsheetService MockSpreadsheetService) writeCell(sheet string, row int, column int, text string) {
	mockSpreadsheetService.Called(sheet, row, column, text)
	return
}

func (mockSpreadsheetService MockSpreadsheetService) readRange(sheet string, rangeToRead string) *sheets.ValueRange {
	args := mockSpreadsheetService.Called(sheet, rangeToRead)
	valueRange := args.Get(0).(sheets.ValueRange)
	return &valueRange
}

type MockSlackService struct {
	mock.Mock
}

func (mockSlackService MockSlackService) PostMessage(channel string, message string) (string, string, error) {
	args := mockSlackService.Called([]string{channel, message})
	return args.String(0), args.String(1), args.Error(2)
}

func (mockSlackService MockSlackService) UpdateMessage(channel, timestamp, text string) (string, string, string, error) {
	args := mockSlackService.Called(channel, timestamp, text)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (mockSlackService MockSlackService) GetUserInfo(user string) (*slack.User, error) {
	args := mockSlackService.Called(user)
	slackUser := args.Get(0).(slack.User)
	return &slackUser, args.Error(1)
}

func (mockSlackService MockSlackService) GetReactions(item slack.ItemRef, params slack.GetReactionsParameters) ([]slack.ItemReaction, error) {
	args := mockSlackService.Called(item, params)
	return args.Get(0).([]slack.ItemReaction), args.Error(1)
}

type MockMessageBuilder struct {
	mock.Mock
}

func (messageBuilder MockMessageBuilder) create(row []interface{}) bytes.Buffer {
	args := messageBuilder.Called(row)
	return args.Get(0).(bytes.Buffer)
}
