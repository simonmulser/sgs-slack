package main

import (
	"bytes"
	"errors"
	"testing"
	"./testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessNew(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"NEW", "05.06.1991 20:04"})
	mockMessageBuilder.On("create", row).Return(createBuffer())
	mockSlackService.On("PostMessage", mock.MatchedBy(func(s []string) bool { return true })).Return("channelID", "timestamp", nil)
	mockSpreadsheetService.On("WriteCell", "topicSheet", 0, main.config.StatusColumn, "POSTED").Return()
	mockSpreadsheetService.On("WriteCell", "topicSheet", 0, main.config.ChannelIDColumn, "channelID").Return()
	mockSpreadsheetService.On("WriteCell", "topicSheet", 0, main.config.TimestampColumn, "timestamp").Return()

	error := eventService.processNew(row, topicConfig, 0)

	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessNewWrongDate(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()
	topicConfig, _ := createTopicConfig()

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"NEW", "05.06.1991"})

	error := eventService.processNew(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}

func TestProcessNewErrorPosting(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"NEW", "05.06.1991 20:04"})
	mockMessageBuilder.On("create", row).Return(createBuffer())
	mockSlackService.On("PostMessage", mock.MatchedBy(func(s []string) bool { return true })).Return("test1", "test2", errors.New("errorFromMock"))

	error := eventService.processNew(row, topicConfig, 0)

	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessPostedWrongDate(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()
	topicConfig, _ := createTopicConfig()

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"POSTED", "05.05.1991 20:20", "topicChannel", "Timestamp", "05.05.1991"})

	error := eventService.processPosted(row, topicConfig, 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "cannot parse")
}

func TestProcessPosted(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	mockSlackService.On("UpdateMessage", "topicChannel", "Timestamp", "~create~").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("WriteCell", "topicSheet", 0, main.config.StatusColumn, "OVER").Return()

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"POSTED", "05.05.1991 20:20", "topicChannel", "Timestamp", "05.05.1991 20:20"})
	mockMessageBuilder.On("create", row).Return(createBuffer())

	error := eventService.processPosted(row, topicConfig, 0)
	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessPostedWithErrorWhileUpdating(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"POSTED", "05.05.1991 20:20", "topicChannel", "Timestamp", "05.05.1991 20:20"})
	mockMessageBuilder.On("create", row).Return(createBuffer())
	mockSlackService.On("UpdateMessage", "topicChannel", "Timestamp", "~create~").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := eventService.processPosted(row, topicConfig, 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessUpdate(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService
	mockSpreadsheetService := new(MockSpreadsheetService)
	main.ISpreadsheetService = mockSpreadsheetService

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"UPDATE", "05.05.1991 20:20", "topicChannel", "Timestamp"})
	mockMessageBuilder.On("create", row).Return(createBuffer())
	mockSlackService.On("UpdateMessage", "topicChannel", "Timestamp", "create").Return("nil", "nil", "nil", nil)
	mockSpreadsheetService.On("WriteCell", "topicSheet", 0, main.config.StatusColumn, "POSTED").Return()

	error := eventService.processUpdate(row, topicConfig, 0)
	assert.Nil(t, error)
	mockMessageBuilder.AssertExpectations(t)
}

func TestProcessUpdateErrorWhileUpdating(t *testing.T) {
	main := Main{}
	main.config = testutils.CreateConfig()

	topicConfig, mockMessageBuilder := createTopicConfig()
	main.IMessageBuilder = mockMessageBuilder
	mockSlackService := new(MockSlackService)
	main.ISlackService = mockSlackService

	eventService := newEventService(&main)

	row := testutils.CreateRow([]string{"UPDATE", "05.05.1991 20:20", "topicChannel", "Timestamp"})
	mockMessageBuilder.On("create", row).Return(createBuffer())
	mockSlackService.On("UpdateMessage", "topicChannel", "Timestamp", "create").Return("nil", "nil", "nil", errors.New("errorFromMock"))

	error := eventService.processUpdate(row, topicConfig, 0)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "errorFromMock")
	mockMessageBuilder.AssertExpectations(t)
}

func createBuffer() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("create")
	return buffer
}

func createTopicConfig() (topicConfig, *MockMessageBuilder) {
	messageBuilder := new(MockMessageBuilder)
	config := topicConfig{
		sheet:           "topicSheet",
		channel:         "topicChannel",
		IMessageBuilder: messageBuilder,
	}
	return config, messageBuilder
}
