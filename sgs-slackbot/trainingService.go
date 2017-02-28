package main

import (
	"time"

	"github.com/golang/glog"
	"github.com/nlopes/slack"
)

type ITopicCommand interface {
	execute(row []interface{}, topic topicConfig, rowNumber int) error
}

// TrainingService processes the trainings
type TrainingService struct {
	main *Main
}

func newTrainingService(main *Main) *TrainingService {
	trainingService := new(TrainingService)
	trainingService.main = main

	return trainingService
}

func (trainingService TrainingService) process() {
	trainingService.postTrainings()
	trainingService.selectResponsibleForUtensils()
	trainingService.strikeTroughOldTrainings()
}

func (trainingService TrainingService) postTrainings() {
	rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TrainingSheet, "A2:G")

	if len(rows.Values) > 0 {
		i := 2
		for _, row := range rows.Values {
			postingDate, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.PostingDateColumn].(string))
			if error != nil {
				glog.Fatalf("Unable to parse date. %v", error)
			}

			if row[trainingService.main.config.ChannelIDColumn] == "FALSE" && timeNow().After(postingDate) {
				message := trainingService.main.messageBuilder.createTrainingPost(row)
				channelID, timestamp, error := trainingService.main.slackService.postMessage(trainingService.main.config.TrainingChannel, message.String())
				if error != nil {
					glog.Fatalf("Unable to post message. %v", error)
				}

				trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TrainingSheet, i, trainingService.main.config.ChannelIDColumn, channelID)
				trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TrainingSheet, i, trainingService.main.config.TimestampColumn, timestamp)
				glog.Info("posted Training and updated sheet")
			}

			i++
		}
	} else {
		glog.Info("No data found.")
	}
}

func (trainingService TrainingService) selectResponsibleForUtensils() {
	rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TrainingSheet, "A2:G")

	if len(rows.Values) > 0 {
		i := 2
		for _, row := range rows.Values {
			date, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.DateColumn].(string))
			if error != nil {
				glog.Fatalf("Unable to parse date. %v", error)
			}
			date = date.Add(-8 * 60 * time.Minute)

			if row[trainingService.main.config.ChannelIDColumn] != "FALSE" && row[trainingService.main.config.TrainingUtensilsColumn] == "FALSE" && timeNow().After(date) {
				reactions, error := trainingService.main.slackService.slack.GetReactions(
					slack.ItemRef{Channel: row[trainingService.main.config.ChannelIDColumn].(string), Timestamp: row[trainingService.main.config.TimestampColumn].(string)},
					slack.GetReactionsParameters{})
				if error != nil {
					glog.Fatalf("Unable to get reactions. %v", error)
				}

				params := trainingService.main.messageBuilder.createTrainingParams(reactions)
				message := trainingService.main.messageBuilder.createTrainingMgmtPost(row, params)
				trainingService.main.slackService.postMessage(trainingService.main.config.TrainingMgmtChannel, message.String())
				trainingService.main.slackService.postMessage("@"+params.ResponsibleTrainingUtensils, trainingService.main.config.TrainingUtensilsResponsibleText)
				trainingService.main.spreadsheetService.writeCell(trainingService.main.config.TrainingSheet, i, trainingService.main.config.TrainingUtensilsColumn, "TRUE")
				glog.Info("selected responsible person and updated sheets")
			}
			i++
		}
	} else {
		glog.Info("No data found.")
	}
}

func (trainingService TrainingService) strikeTroughOldTrainings() {
	rows := trainingService.main.spreadsheetService.readRange(trainingService.main.config.TrainingSheet, "A2:G")

	if len(rows.Values) > 0 {
		i := 2
		for _, row := range rows.Values {
			if row[trainingService.main.config.ChannelIDColumn] != "FALSE" && row[trainingService.main.config.ChannelIDColumn] != "TRUE" {
				date, error := time.Parse("02.01.2006 15:04", row[trainingService.main.config.DateColumn].(string))
				if error != nil {
					glog.Fatalf("Unable to parse date. %v", error)
				}

				date = date.Add(12 * time.Hour)
				if timeNow().After(date) {
					message := trainingService.main.messageBuilder.createTrainingPost(row)
					trainingService.main.slackService.slack.UpdateMessage(row[trainingService.main.config.ChannelIDColumn].(string), row[trainingService.main.config.TimestampColumn].(string),
						"~"+message.String()+"~")
					if error != nil {
						glog.Fatalf("Unable to post massage. %v", error)
					}
					glog.Info("updated Training")
				}
			}

			i++
		}
	} else {
		glog.Info("No data found.")
	}
}
