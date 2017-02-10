package main

import (
	"time"

	"github.com/golang/glog"
)

// GameService to process games
type GameService struct {
	main               *Main
	spreadsheetService *SpreadsheetService
	teams              []teamConfig
}

type teamConfig struct {
	sheet   string
	channel string
}

func newGameService(main *Main) *GameService {
	GameService := new(GameService)
	GameService.main = main

	sgs07 := teamConfig{main.config.Games07Sheet, main.config.Games07Channel}
	sgs16 := teamConfig{main.config.Games16Sheet, main.config.Games16Channel}
	GameService.teams = []teamConfig{sgs07, sgs16}

	return GameService
}

func (gameService GameService) process() {
	gameService.postGames()

	gameService.strikeTroughOldGames()
}

func (gameService GameService) postGames() {
	for _, team := range gameService.teams {
		rows := gameService.main.spreadsheetService.readRange(team.sheet, "A2:K")

		if len(rows.Values) > 0 {
			i := 2
			for _, row := range rows.Values {
				postingDate, error := time.Parse("02.01.2006 15:04", row[gameService.main.config.GamePostingDateColumn].(string))
				if error != nil {
					glog.Fatalf("Unable to parse date. %v", error)
				}

				if row[gameService.main.config.GameChannelIDColumn] == "FALSE" && timeNow().After(postingDate) {
					message := gameService.main.messageBuilder.createGamePost(row)
					channelID, timestamp, error := gameService.main.slackService.postMessage(team.channel, message.String())
					if error != nil {
						glog.Fatalf("Unable to post massage. %v", error)
					}

					gameService.main.spreadsheetService.writeCell(team.sheet, i, gameService.main.config.GameChannelIDColumn, channelID)
					gameService.main.spreadsheetService.writeCell(team.sheet, i, gameService.main.config.GameTimestampColumn, timestamp)
					glog.Info("posted Game and updated sheet")
				}

				i++
			}
		} else {
			glog.Info("No data found.")
		}
	}
}

func (gameService GameService) strikeTroughOldGames() {
	for _, team := range gameService.teams {
		rows := gameService.main.spreadsheetService.readRange(team.sheet, "A2:K")

		if len(rows.Values) > 0 {
			i := 2
			for _, row := range rows.Values {
				if row[gameService.main.config.GameChannelIDColumn] != "FALSE" && row[gameService.main.config.GameChannelIDColumn] != "TRUE" {
					date, error := time.Parse("02.01.2006 15:04", row[gameService.main.config.GameDateColumn].(string))
					if error != nil {
						glog.Fatalf("Unable to parse date. %v", error)
					}

					date = date.Add(12 * time.Hour)
					if timeNow().After(date) {
						message := gameService.main.messageBuilder.createGamePost(row)
						gameService.main.slackService.slack.UpdateMessage(row[gameService.main.config.GameChannelIDColumn].(string), row[gameService.main.config.GameTimestampColumn].(string),
							"~"+message.String()+"~")
						if error != nil {
							glog.Fatalf("Unable to post massage. %v", error)
						}
						glog.Info("updated Game")
					}
				}

				i++
			}
		} else {
			glog.Info("No data found.")
		}
	}
}
