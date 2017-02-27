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

	eventService := newEventService(&main)

	row := createRow([]string{"NEW", "05.06.1991 20:04"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createEventPost").Return("channelID", "timestamp", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.StatusColumn, "POSTED").Return()
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.ChannelIDColumn, "channelID").Return()
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.TimestampColumn, "timestamp").Return()

	error := eventService.processNew(row, createTeamConfig(), 0)

	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessNewWrongDate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	eventService := newEventService(&main)

	row := createRow([]string{"NEW", "05.06.1991"})

	error := eventService.processNew(row, createTeamConfig(), 0)
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

	eventService := newEventService(&main)

	row := createRow([]string{"NEW", "05.06.1991 20:04"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())
	mockSlackService.On("postMessage", "teamChannel", "createEventPost").Return("test1", "test2", errors.New("errorFromMock"))

	error := eventService.processNew(row, createTeamConfig(), 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessPostedWrongDate(t *testing.T) {
	main := Main{}
	main.config = createConfig()

	eventService := newEventService(&main)

	row := createRow([]string{"POSTED", "05.05.1991"})

	error := eventService.processPosted(row, createTeamConfig(), 0)
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

	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "~createEventPost~").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.StatusColumn, "OVER").Return()

	eventService := newEventService(&main)

	row := createRow([]string{"POSTED", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())

	error := eventService.processPosted(row, createTeamConfig(), 0)
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

	eventService := newEventService(&main)

	row := createRow([]string{"POSTED", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "~createEventPost~").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := eventService.processPosted(row, createTeamConfig(), 0)
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

	eventService := newEventService(&main)

	row := createRow([]string{"UPDATE", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "createEventPost").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("writeCell", "teamSheet", 0, main.config.StatusColumn, "POSTED").Return()

	error := eventService.processUpdate(row, createTeamConfig(), 0)
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

	eventService := newEventService(&main)

	row := createRow([]string{"UPDATE", "05.05.1991 20:20", "teamChannel", "Timestamp"})
	mockMessageBuilder.On("createEventPost", row).Return(createBuffer())
	mockSlackService.On("updateMessage", "teamChannel", "Timestamp", "createEventPost").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := eventService.processUpdate(row, createTeamConfig(), 0)
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
		StatusColumn:      0,
		PostingDateColumn: 1,
		ChannelIDColumn:   2,
		TimestampColumn:   3,
		DateColumn:        1,
	}
	return &config
}

func createBuffer() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("createEventPost")
	return buffer
}

func createTeamConfig() teamConfig {
	config := teamConfig{
		sheet:   "teamSheet",
		channel: "teamChannel",
	}
	return config
}
