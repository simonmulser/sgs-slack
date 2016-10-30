package main

import (
  "github.com/golang/glog"
  "time"
  "os"
  "flag"
  
  "github.com/nlopes/slack"
  "google.golang.org/api/sheets/v4"
)

type Main struct{
  config *Config
  slackClient *slack.Client
  service *sheets.Service
}

func main() {
  flag.Parse()
  flag.Lookup("logtostderr").Value.Set("true")
  glog.Info("Program started")

  instance := Main{}
  instance.run()

  glog.Info("Program terminated")
}

func (main Main) run() {
  env := "development"
  if (len(os.Args) > 1) {
    env = os.Args[1]
  }
  main.config = Read(env)

  main.service = New();

  main.slackClient = main.createSlackClient(main.config.SLACK_KEY)

  var rows = main.readRange(main.config.TRAINING_SHEET ,"A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
        message := main.createTrainingPost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())
        if error != nil {
          glog.Fatalf("Unable to post message. %v", error)
        }

        main.writeCell(main.config.TRAINING_SHEET, i, main.config.CHANNEL_ID_COLUMN, channelId)
        main.writeCell(main.config.TRAINING_SHEET, i, main.config.TIMESTAMP_COLUMN, timestamp)
      }

      i++
      }
    } else {
      glog.Info("No data found.")
    }

  rows = main.readRange(main.config.TRAINING_SHEET, "A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      date, error := time.Parse("02.01.2006 15:04", row[main.config.DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }
      date = date.Add(-8 * 60 * time.Minute)

      if(row[main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[main.config.BALLS_COLUMN] == "FALSE" && timeNow().After(date)){
          reactions, error := main.slackClient.GetReactions(
            slack.ItemRef{Channel: row[main.config.CHANNEL_ID_COLUMN].(string), Timestamp: row[main.config.TIMESTAMP_COLUMN].(string)},
            slack.GetReactionsParameters{})
          if error != nil {
            glog.Fatalf("Unable to get reactions. %v", error)
          }

          params := main.createTrainingParams(reactions)
          message := main.createTrainingMgmtPost(row, params)
          main.postMessage(main.config.TRAINING_MGMT_CHANNEL, message.String())
          main.postMessage("@" + params.Responsible_balls, main.config.BALLS_RESPONSIBLE_TEXT)
          main.writeCell(main.config.TRAINING_SHEET, i, main.config.BALLS_COLUMN, "TRUE")
          }
      i++
      }
    } else {
      glog.Info("No data found.")
    }

if false {
  rows = main.readRange(main.config.GAMES_07_SHEET, "A2:K")

    if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.GAME_POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.GAME_CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
        message := main.createGamePost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())
        if error != nil {
          glog.Fatalf("Unable to post massage. %v", error)
        }

        main.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_CHANNEL_ID_COLUMN, channelId)
        main.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_TIMESTAMP_COLUMN, timestamp)
      }

      i++
      }
    } else {
      glog.Info("No data found.")
    }
  }
}
