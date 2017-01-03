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
