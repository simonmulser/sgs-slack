package main

import (
  "time"

  "github.com/golang/glog"
  "github.com/nlopes/slack"
)

type TrainingService struct{
  main *Main
}

func NewTrainingService(main *Main) *TrainingService {
  trainingService := new(TrainingService)
  trainingService.main = main

  return trainingService;
}

func (trainingService TrainingService) process() {
  trainingService.postTrainings()
  trainingService.selectResponsibleForUtensils()
  trainingService.strikeTroughOldTrainings()
}

func (trainingService TrainingService) postTrainings() {
rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TRAINING_SHEET ,"A2:G")

if len(rows.Values) > 0 {
  i := 2
  for _, row := range rows.Values {
    postingDate, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.POSTING_DATE_COLUMN].(string))
    if error != nil {
      glog.Fatalf("Unable to parse date. %v", error)
    }

    if(row[trainingService.main.config.CHANNEL_ID_COLUMN] == "FALSE" && timeNow().After(postingDate)){
      message := trainingService.main.messageBuilder.createTrainingPost(row)
      channelId, timestamp, error := trainingService.main.slackService.postMessage(trainingService.main.config.TRAINING_CHANNEL, message.String())
      if error != nil {
        glog.Fatalf("Unable to post message. %v", error)
      }

      trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TRAINING_SHEET, i, trainingService.main.config.CHANNEL_ID_COLUMN, channelId)
      trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TRAINING_SHEET, i, trainingService.main.config.TIMESTAMP_COLUMN, timestamp)
      glog.Info("posted training and updated sheet")
    }

    i++
    }
  } else {
    glog.Info("No data found.")
  }
}

func (trainingService TrainingService) selectResponsibleForUtensils() {
rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TRAINING_SHEET, "A2:G")

if len(rows.Values) > 0 {
  i := 2
  for _, row := range rows.Values {
    date, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.DATE_COLUMN].(string))
    if error != nil {
      glog.Fatalf("Unable to parse date. %v", error)
    }
    date = date.Add(-8 * 60 * time.Minute)

    if(row[trainingService.main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[trainingService.main.config.TRAINING_UTENSILS_COLUMN] == "FALSE" && timeNow().After(date)){
        reactions, error := trainingService.main.slackService.slack.GetReactions(
          slack.ItemRef{Channel: row[trainingService.main.config.CHANNEL_ID_COLUMN].(string), Timestamp: row[trainingService.main.config.TIMESTAMP_COLUMN].(string)},
          slack.GetReactionsParameters{})
        if error != nil {
          glog.Fatalf("Unable to get reactions. %v", error)
        }

        params := trainingService.main.messageBuilder.createTrainingParams(reactions)
        message := trainingService.main.messageBuilder.createTrainingMgmtPost(row, params)
        trainingService.main.slackService.postMessage(trainingService.main.config.TRAINING_MGMT_CHANNEL, message.String())
        trainingService.main.slackService.postMessage("@" + params.Responsible_training_utensils, trainingService.main.config.TRAINING_UTENSILS_RESPONSIBLE_TEXT)
        trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TRAINING_SHEET, i, trainingService.main.config.TRAINING_UTENSILS_COLUMN, "TRUE")
        glog.Info("selected responsible person and updated sheets")
        }
    i++
    }
  } else {
    glog.Info("No data found.")
  }
}

func (trainingService TrainingService) strikeTroughOldTrainings() {
  rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TRAINING_SHEET, "A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      if(row[trainingService.main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[trainingService.main.config.CHANNEL_ID_COLUMN] != "TRUE"){
        date, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.DATE_COLUMN].(string))
        if error != nil {
          glog.Fatalf("Unable to parse date. %v", error)
        }

        date = date.Add(24 * time.Hour)
        if(timeNow().After(date)) {
          message := trainingService.main.messageBuilder.createTrainingPost(row)
          trainingService.main.slackService.slack.UpdateMessage(row[trainingService.main.config.CHANNEL_ID_COLUMN].(string), row[trainingService.main.config.TIMESTAMP_COLUMN].(string),
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
