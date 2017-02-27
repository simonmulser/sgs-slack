package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
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
	var buffer bytes.Buffer
	return buffer
}

func (messageBuilder MockMessageBuilder) createTrainingParams(reactions []slack.ItemReaction) trainingParameters {
	var params trainingParameters
	return params
}

func TestToProcess(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	gameService := newGameService(&main)

	row := createRow([]string{"PLANEND", "05.06.1991 20:04", "FALSE"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createGamePost").Return("channelID", "timestamp", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameChannelIDColumn, "channelID").Return()
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameTimestampColumn, "timestamp").Return()

	error := gameService.toProcess(row, createTeamConfig(), 0)

	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestToProcessWrongDate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	gameService := newGameService(&main)

	row := createRow([]string{"PLANEND", "05.06.1991", "FALSE"})

	error := gameService.toProcess(row, createTeamConfig(), 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}

func TestToProcessErrorPosting(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	gameService := newGameService(&main)

	row := createRow([]string{"PLANEND", "05.06.1991 20:04", "FALSE"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createGamePost").Return("test1", "test2", errors.New("errorFromMock"))

	error := gameService.toProcess(row, createTeamConfig(), 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func createRow(data []string) []interface{} {
	row := make([]interface{}, len(data))
	for i, s := range data {
		row[i] = s
	}

	return row
}

func createConfig() *Config {
	config := Config{
		GameStatusColumn:      0,
		GamePostingDateColumn: 1,
		GameChannelIDColumn:   2,
		GameTimestampColumn:   3,
	}
	return &config
}

func createBuffer() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("createGamePost")
	return buffer
}

func createTeamConfig() teamConfig {
	config := teamConfig{
		sheet:   "teamSheet",
		channel: "teamChannel",
	}
	return config
}