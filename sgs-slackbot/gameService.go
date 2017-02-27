package main

import (
	"time"

	"github.com/golang/glog"
)

// GameService to process games
type GameService struct {
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

func newGameService(main *Main) *GameService {
	gameService := new(GameService)
	gameService.config = main.config
	gameService.IMessageBuilder = main.IMessageBuilder
	gameService.ISlackService = main.ISlackService
	gameService.ISpreadsheetService = main.ISpreadsheetService

	sgs07 := teamConfig{main.config.Games07Sheet, main.config.Games07Channel}
	sgs16 := teamConfig{main.config.Games16Sheet, main.config.Games16Channel}
	gameService.teams = []teamConfig{sgs07, sgs16}

	return gameService
}

func (gameService GameService) process() {
	for _, team := range gameService.teams {
		rows := gameService.ISpreadsheetService.readRange(team.sheet, "A2:L")
		if len(rows.Values) > 0 {
			i := 2
			for _, row := range rows.Values {
				var error error
				switch row[gameService.config.StatusColumn] {
				case "NEW":
					error = gameService.processNew(row, team, i)
				case "POSTED":
					error = gameService.processPosted(row, team, i)
				case "UPDATE":
					error = gameService.processUpdate(row, team, i)
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

func (gameService GameService) processNew(row []interface{}, team teamConfig, rowNumber int) error {
	postingDate, error := time.Parse("02.01.2006 15:04", row[gameService.config.PostingDateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	if timeNow().After(postingDate) {
		message := gameService.IMessageBuilder.createGamePost(row)
		channelID, timestamp, error := gameService.ISlackService.postMessage(team.channel, message.String())
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		gameService.ISpreadsheetService.writeCell(team.sheet, rowNumber, gameService.config.StatusColumn, "POSTED")
		gameService.ISpreadsheetService.writeCell(team.sheet, rowNumber, gameService.config.ChannelIDColumn, channelID)
		gameService.ISpreadsheetService.writeCell(team.sheet, rowNumber, gameService.config.TimestampColumn, timestamp)
		glog.Info("posted Game and updated sheet")
	}
	return nil
}

func (gameService GameService) processPosted(row []interface{}, team teamConfig, rowNumber int) error {
	date, error := time.Parse("02.01.2006 15:04", row[gameService.config.DateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}

	date = date.Add(6 * time.Hour)

	if timeNow().After(date) {
		message := gameService.IMessageBuilder.createGamePost(row)
		_, _, _, error := gameService.ISlackService.updateMessage(row[gameService.config.ChannelIDColumn].(string), row[gameService.config.TimestampColumn].(string),
			"~"+message.String()+"~")
		if error != nil {
			glog.Warningf("Unable to post massage. %v", error)
			return error
		}

		gameService.ISpreadsheetService.writeCell(team.sheet, rowNumber, gameService.config.StatusColumn, "OVER")
		glog.Info("updated Game")
	}
	return nil
}

func (gameService GameService) processUpdate(row []interface{}, team teamConfig, rowNumber int) error {
	message := gameService.IMessageBuilder.createGamePost(row)
	_, _, _, error := gameService.ISlackService.updateMessage(row[gameService.config.ChannelIDColumn].(string), row[gameService.config.TimestampColumn].(string), message.String())

	if error != nil {
		glog.Warningf("Unable to post massage. %v", error)
		return error
	}

	gameService.ISpreadsheetService.writeCell(team.sheet, rowNumber, gameService.config.StatusColumn, "POSTED")
	glog.Info("updated Game")
	return nil
}
