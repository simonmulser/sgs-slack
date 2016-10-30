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
  slackService *SlackService
  service *sheets.Service
  messageBuilder *MessageBuilder
}

func main() {
  flag.Parse()
  flag.Lookup("logtostderr").Value.Set("true")
  glog.Info("Program started")

  env := "development"
  if (len(os.Args) > 1) {
    env = os.Args[1]
  }

  instance := Main{}
  instance.config = Read(env)
  instance.slackService = NewSlackService(instance.config.SLACK_KEY)
  instance.messageBuilder = NewMessageBuilder(instance.config, instance.slackService)
  instance.service = New();

  instance.run()

  glog.Info("Program terminated")
}

func (main Main) run() {
  var rows = main.readRange(main.config.TRAINING_SHEET ,"A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
        message := main.messageBuilder.createTrainingPost(row)
        channelId, timestamp, error := main.slackService.postMessage(main.config.TRAINING_CHANNEL, message.String())
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

      if(row[main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[main.config.TRAINING_UTENSILS_COLUMN] == "FALSE" && timeNow().After(date)){
          reactions, error := main.slackService.slack.GetReactions(
            slack.ItemRef{Channel: row[main.config.CHANNEL_ID_COLUMN].(string), Timestamp: row[main.config.TIMESTAMP_COLUMN].(string)},
            slack.GetReactionsParameters{})
          if error != nil {
            glog.Fatalf("Unable to get reactions. %v", error)
          }

          params := main.messageBuilder.createTrainingParams(reactions)
          message := main.messageBuilder.createTrainingMgmtPost(row, params)
          main.slackService.postMessage(main.config.TRAINING_MGMT_CHANNEL, message.String())
          main.slackService.postMessage("@" + params.Responsible_training_utensils, main.config.TRAINING_UTENSILS_RESPONSIBLE_TEXT)
          main.writeCell(main.config.TRAINING_SHEET, i, main.config.TRAINING_UTENSILS_COLUMN, "TRUE")
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
        message := main.messageBuilder.createGamePost(row)
        channelId, timestamp, error := main.slackService.postMessage(main.config.TRAINING_CHANNEL, message.String())
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
