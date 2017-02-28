package main

import (
	"time"

	"github.com/golang/glog"
)

// EventService to process events
type EventService struct {
	config *Config
	topics []topicConfig
	IMessageBuilder
	ISlackService
	ISpreadsheetService
}

type topicConfig struct {
	sheet   string
	channel string
	IMessageBuilder
	ITopicCommand
}

func newEventService(main *Main) *EventService {
	eventService := new(EventService)
	eventService.config = main.config
	eventService.ISlackService = main.ISlackService
	eventService.ISpreadsheetService = main.ISpreadsheetService

	sgs07 := topicConfig{main.config.Games07Sheet, main.config.Games07Channel, newGameMessageBuilder(eventService.config), nil}
	sgs16 := topicConfig{main.config.Games16Sheet, main.config.Games16Channel, newGameMessageBuilder(eventService.config), nil}
	training := topicConfig{main.config.TrainingSheet, main.config.TrainingChannel, newGameMessageBuilder(eventService.config), nil}
	eventService.topics = []topicConfig{sgs07, sgs16, training}

	return eventService
}

func (eventService EventService) process() {
	for _, topic := range eventService.topics {
		rows := eventService.ISpreadsheetService.readRange(topic.sheet, "A2:L")
		if len(rows.Values) > 0 {
			i := 2
			for _, row := range rows.Values {
				var error error
				switch row[eventService.config.StatusColumn] {
				case "NEW":
					error = eventService.processNew(row, topic, i)
				case "POSTED":
					error = eventService.processPosted(row, topic, i)
				case "UPDATE":
					error = eventService.processUpdate(row, topic, i)
				case "OVER":
				default:
					glog.Warningf("Status not recognized")
				}

				if error != nil {
					glog.Warningf("Could not process row %v", error)
				}

				if topic.ITopicCommand != nil {
					topic.ITopicCommand.execute(row, topic, i)
				}

				i++
			}
		} else {
			glog.Info("No data found.")
		}
	}
}

func (eventService EventService) processNew(row []interface{}, topic topicConfig, rowNumber int) error {
	postingDate, error := time.Parse("02.01.2006 15:04", row[eventService.config.PostingDateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	if timeNow().After(postingDate) {
		message := topic.IMessageBuilder.create(row)
		channelID, timestamp, error := eventService.ISlackService.postMessage(topic.channel, message.String())
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		eventService.ISpreadsheetService.writeCell(topic.sheet, rowNumber, eventService.config.StatusColumn, "POSTED")
		eventService.ISpreadsheetService.writeCell(topic.sheet, rowNumber, eventService.config.ChannelIDColumn, channelID)
		eventService.ISpreadsheetService.writeCell(topic.sheet, rowNumber, eventService.config.TimestampColumn, timestamp)
		glog.Info("posted event and updated sheet")
	}
	return nil
}

func (eventService EventService) processPosted(row []interface{}, topic topicConfig, rowNumber int) error {
	date, error := time.Parse("02.01.2006 15:04", row[eventService.config.DateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	date = date.Add(6 * time.Hour)

	if timeNow().After(date) {
		message := topic.IMessageBuilder.create(row)
		_, _, _, error := eventService.ISlackService.updateMessage(row[eventService.config.ChannelIDColumn].(string), row[eventService.config.TimestampColumn].(string),
			"~"+message.String()+"~")
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		eventService.ISpreadsheetService.writeCell(topic.sheet, rowNumber, eventService.config.StatusColumn, "OVER")
		glog.Info("updated event")
	}
	return nil
}

func (eventService EventService) processUpdate(row []interface{}, topic topicConfig, rowNumber int) error {
	message := topic.IMessageBuilder.create(row)
	_, _, _, error := eventService.ISlackService.updateMessage(row[eventService.config.ChannelIDColumn].(string), row[eventService.config.TimestampColumn].(string), message.String())

	if error != nil {
		glog.Warningf("Unable to post massage. %v", error)
		return error
	}

	eventService.ISpreadsheetService.writeCell(topic.sheet, rowNumber, eventService.config.StatusColumn, "POSTED")
	glog.Info("updated event")
	return nil
}
