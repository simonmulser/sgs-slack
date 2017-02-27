package main

import (
	"time"

	"github.com/golang/glog"
)

// EventService to process events
type EventService struct {
	config *Config
	teams  []teamConfig
	IMessageBuilder
	ISlackService
	ISpreadsheetService
}

type teamConfig struct {
	sheet   string
	channel string
}

func newEventService(main *Main) *EventService {
	eventService := new(EventService)
	eventService.config = main.config
	eventService.IMessageBuilder = main.IMessageBuilder
	eventService.ISlackService = main.ISlackService
	eventService.ISpreadsheetService = main.ISpreadsheetService

	sgs07 := teamConfig{main.config.Games07Sheet, main.config.Games07Channel}
	sgs16 := teamConfig{main.config.Games16Sheet, main.config.Games16Channel}
	eventService.teams = []teamConfig{sgs07, sgs16}

	return eventService
}

func (eventService EventService) process() {
	for _, team := range eventService.teams {
		rows := eventService.ISpreadsheetService.readRange(team.sheet, "A2:L")
		if len(rows.Values) > 0 {
			i := 2
			for _, row := range rows.Values {
				var error error
				switch row[eventService.config.StatusColumn] {
				case "NEW":
					error = eventService.processNew(row, team, i)
				case "POSTED":
					error = eventService.processPosted(row, team, i)
				case "UPDATE":
					error = eventService.processUpdate(row, team, i)
				case "OVER":
				default:
					glog.Warningf("Status not recognized")
				}

				if error != nil {
					glog.Warningf("Could not process row %v", error)
				}

				i++
			}
		} else {
			glog.Info("No data found.")
		}
	}
}

func (eventService EventService) processNew(row []interface{}, team teamConfig, rowNumber int) error {
	postingDate, error := time.Parse("02.01.2006 15:04", row[eventService.config.PostingDateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	if timeNow().After(postingDate) {
		message := eventService.IMessageBuilder.createEventPost(row)
		channelID, timestamp, error := eventService.ISlackService.postMessage(team.channel, message.String())
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		eventService.ISpreadsheetService.writeCell(team.sheet, rowNumber, eventService.config.StatusColumn, "POSTED")
		eventService.ISpreadsheetService.writeCell(team.sheet, rowNumber, eventService.config.ChannelIDColumn, channelID)
		eventService.ISpreadsheetService.writeCell(team.sheet, rowNumber, eventService.config.TimestampColumn, timestamp)
		glog.Info("posted event and updated sheet")
	}
	return nil
}

func (eventService EventService) processPosted(row []interface{}, team teamConfig, rowNumber int) error {
	date, error := time.Parse("02.01.2006 15:04", row[eventService.config.DateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	date = date.Add(6 * time.Hour)

	if timeNow().After(date) {
		message := eventService.IMessageBuilder.createEventPost(row)
		_, _, _, error := eventService.ISlackService.updateMessage(row[eventService.config.ChannelIDColumn].(string), row[eventService.config.TimestampColumn].(string),
			"~"+message.String()+"~")
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		eventService.ISpreadsheetService.writeCell(team.sheet, rowNumber, eventService.config.StatusColumn, "OVER")
		glog.Info("updated event")
	}
	return nil
}

func (eventService EventService) processUpdate(row []interface{}, team teamConfig, rowNumber int) error {
	message := eventService.IMessageBuilder.createEventPost(row)
	_, _, _, error := eventService.ISlackService.updateMessage(row[eventService.config.ChannelIDColumn].(string), row[eventService.config.TimestampColumn].(string), message.String())

	if error != nil {
		glog.Warningf("Unable to post massage. %v", error)
		return error
	}

	eventService.ISpreadsheetService.writeCell(team.sheet, rowNumber, eventService.config.StatusColumn, "POSTED")
	glog.Info("updated event")
	return nil
}
