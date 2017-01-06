package main

import (
  "github.com/golang/glog"
  "time"
)

type GameService struct{
  main *Main
  spreadsheetService *SpreadsheetService
  teams []TeamConfig
}

type TeamConfig struct{
  sheet string
  channel string
}

func NewGameService(main *Main) *GameService {
  gameService := new(GameService)
  gameService.main = main

  sgs07 := TeamConfig{main.config.GAMES_07_SHEET, main.config.GAMES_07_CHANNEL}
  sgs16 := TeamConfig{main.config.GAMES_16_SHEET, main.config.GAMES_16_CHANNEL}
  gameService.teams = []TeamConfig{sgs07, sgs16}

  return gameService;
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
        postingDate, error := time.Parse("02.01.2006 15:04", row[gameService.main.config.GAME_POSTING_DATE_COLUMN].(string))
        if error != nil {
          glog.Fatalf("Unable to parse date. %v", error)
        }

        if(row[gameService.main.config.GAME_CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
          message := gameService.main.messageBuilder.createGamePost(row)
          channelId, timestamp, error := gameService.main.slackService.postMessage(team.channel, message.String())
          if error != nil {
            glog.Fatalf("Unable to post massage. %v", error)
          }

          gameService.main.spreadsheetService.writeCell(team.sheet, i, gameService.main.config.GAME_CHANNEL_ID_COLUMN, channelId)
          gameService.main.spreadsheetService.writeCell(team.sheet, i, gameService.main.config.GAME_TIMESTAMP_COLUMN, timestamp)
          glog.Info("posted game and updated sheet")
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
        if(row[gameService.main.config.GAME_CHANNEL_ID_COLUMN] != "FALSE" && row[gameService.main.config.GAME_CHANNEL_ID_COLUMN] != "TRUE"){
          date, error := time.Parse("02.01.2006 15:04", row[gameService.main.config.GAME_DATE_COLUMN].(string))
          if error != nil {
            glog.Fatalf("Unable to parse date. %v", error)
          }

          date = date.Add(24 * time.Hour)
          if(timeNow().After(date)) {
            message := gameService.main.messageBuilder.createGamePost(row)
            gameService.main.slackService.slack.UpdateMessage(row[gameService.main.config.GAME_CHANNEL_ID_COLUMN].(string), row[gameService.main.config.GAME_TIMESTAMP_COLUMN].(string),
                "~" + message.String() + "~")
            if error != nil {
              glog.Fatalf("Unable to post massage. %v", error)
            }
            glog.Info("updated game")
          }
        }

        i++
      }
    } else {
      glog.Info("No data found.")
    }
  }
}
