package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessNew(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	gameService := newGameService(&main)

	row := createRow([]string{"NEW", "05.06.1991 20:04", "FALSE"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createGamePost").Return("channelID", "timestamp", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameStatusColumn, "POSTED").Return()
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameChannelIDColumn, "channelID").Return()
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameTimestampColumn, "timestamp").Return()

	error := gameService.processNew(row, createTeamConfig(), 0)

	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessNewWrongDate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	gameService := newGameService(&main)

	row := createRow([]string{"NEW", "05.06.1991", "FALSE"})

	error := gameService.processNew(row, createTeamConfig(), 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}

func TestProcessNewErrorPosting(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	gameService := newGameService(&main)

	row := createRow([]string{"NEW", "05.06.1991 20:04", "FALSE"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createGamePost").Return("test1", "test2", errors.New("errorFromMock"))

	error := gameService.processNew(row, createTeamConfig(), 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessPostedWrongDate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	gameService := newGameService(&main)

	row := createRow([]string{"POSTED", "05.05.1991"})

	error := gameService.processPosted(row, createTeamConfig(), 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}

func TestProcessPosted(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "~createGamePost~").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameStatusColumn, "OVER").Return()

	gameService := newGameService(&main)

	row := createRow([]string{"POSTED", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())

	error := gameService.processPosted(row, createTeamConfig(), 0)
	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessPostedWithErrorWhileUpdating(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	gameService := newGameService(&main)

	row := createRow([]string{"POSTED", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "~createGamePost~").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := gameService.processPosted(row, createTeamConfig(), 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessUpdate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	gameService := newGameService(&main)

	row := createRow([]string{"UPDATE", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "createGamePost").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.GameStatusColumn, "POSTED").Return()

	error := gameService.processUpdate(row, createTeamConfig(), 0)
	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessUpdateErrorWhileUpdating(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	mockMessageBuilder := new(MockMessageBuilder)
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	gameService := newGameService(&main)

	row := createRow([]string{"UPDATE", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createGamePost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "createGamePost").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := gameService.processUpdate(row, createTeamConfig(), 0)
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
		GameDateColumn:        1,
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
