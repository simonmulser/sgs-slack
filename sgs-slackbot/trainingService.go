package main

import (
  "time"

  "github.com/golang/glog"
  "github.com/nlopes/slack"
)

func processTrainings(main Main) {
  postTraining(main)
  selectResponsibleForUtensils(main)
}

func postTraining(main Main) {
rows := main.spreadsheetService.readRange(main.config.TRAINING_SHEET ,"A2:G")

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

      main.spreadsheetService.writeCell(main.config.TRAINING_SHEET, i, main.config.CHANNEL_ID_COLUMN, channelId)
      main.spreadsheetService.writeCell(main.config.TRAINING_SHEET, i, main.config.TIMESTAMP_COLUMN, timestamp)
    }

    i++
    }
  } else {
    glog.Info("No data found.")
  }
}

func selectResponsibleForUtensils(main Main) {
rows := main.spreadsheetService.readRange(main.config.TRAINING_SHEET, "A2:G")

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
        main.spreadsheetService.writeCell(main.config.TRAINING_SHEET, i, main.config.TRAINING_UTENSILS_COLUMN, "TRUE")
        }
    i++
    }
  } else {
    glog.Info("No data found.")
  }
}
