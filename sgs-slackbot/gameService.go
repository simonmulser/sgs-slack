package main

import (
  "github.com/golang/glog"
  "time"
)

func processGames(main Main) {
  rows := main.spreadsheetService.readRange(main.config.GAMES_07_SHEET, "A2:K")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.GAME_POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.GAME_CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
        message := main.messageBuilder.createGamePost(row)
        channelId, timestamp, error := main.slackService.postMessage(main.config.GAMES_07_CHANNEL, message.String())
        if error != nil {
          glog.Fatalf("Unable to post massage. %v", error)
        }

        main.spreadsheetService.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_CHANNEL_ID_COLUMN, channelId)
        main.spreadsheetService.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_TIMESTAMP_COLUMN, timestamp)
        glog.Info("posted game and updated sheet")
      }

      i++
    }
  } else {
    glog.Info("No data found.")
  }
}

func strikeTroughOldGames(main Main) {
  rows := main.spreadsheetService.readRange(main.config.GAMES_07_SHEET, "A2:K")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      if(row[main.config.GAME_CHANNEL_ID_COLUMN] != "FALSE" && row[main.config.GAME_CHANNEL_ID_COLUMN] != "TRUE"){
        date, error := time.Parse("02.01.2006 15:04", row[main.config.GAME_DATE_COLUMN].(string))
        if error != nil {
          glog.Fatalf("Unable to parse date. %v", error)
        }

        date = date.Add(24 * time.Hour)
        if(timeNow().After(date)) {
          message := main.messageBuilder.createGamePost(row)
          main.slackService.slack.UpdateMessage(row[main.config.GAME_CHANNEL_ID_COLUMN].(string), row[main.config.GAME_TIMESTAMP_COLUMN].(string),
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
