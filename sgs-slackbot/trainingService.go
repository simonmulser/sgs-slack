package main

import (
	"bytes"
	"time"
	"./config"
	"./google"
	"./slack"
	
	"github.com/golang/glog"
	nlopesslack "github.com/nlopes/slack"
)

type ITopicCommand interface {
	execute(row []interface{}, topic topicConfig, rowNumber int) error
}

// TrainingService processes the trainings
type TrainingService struct {
	config *config.Config
	slack.ISlackService
	google.ISpreadsheetService
	ITrainingParamsService
}

func newTrainingService(config *config.Config, slackService slack.ISlackService, spreadsheetService google.ISpreadsheetService, trainingParamsService ITrainingParamsService) *TrainingService {
	trainingService := new(TrainingService)
	trainingService.config = config
	trainingService.ISlackService = slackService
	trainingService.ISpreadsheetService = spreadsheetService
	trainingService.ITrainingParamsService = trainingParamsService

	return trainingService
}

func (trainingService TrainingService) execute(row []interface{}, topic topicConfig, rowNumber int) error {
	date, error := time.Parse("02/01/2006 15:04", row[trainingService.config.DateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}
	date = date.Add(-8 * 60 * time.Minute)

	if row[trainingService.config.StatusColumn] == "POSTED" && row[trainingService.config.TrainingUtensilsColumn] != "POSTED" && timeNow().After(date) {
		reactions, error := trainingService.ISlackService.GetReactions(
			nlopesslack.ItemRef{Channel: row[trainingService.config.ChannelIDColumn].(string), Timestamp: row[trainingService.config.TimestampColumn].(string)},
			nlopesslack.GetReactionsParameters{})
		if error != nil {
			glog.Warningf("Unable to get reactions. %v", error)
			return error
		}
		glog.Infof("Got reactions from post (channel=%s, timestamp=%s)", row[trainingService.config.ChannelIDColumn], row[trainingService.config.TimestampColumn])

		params := trainingService.ITrainingParamsService.create(reactions)
		message := trainingService.createTrainingMgmtPost(row, params)

		trainingService.ISlackService.PostMessage(trainingService.config.TrainingMgmtChannel, message.String())
		trainingService.ISlackService.PostMessage("@"+params.ResponsibleTrainingUtensils, trainingService.config.TrainingUtensilsResponsibleText)
		glog.Infof("Informed responsible person=%s and posted into channel=%s", params.ResponsibleTrainingUtensils, trainingService.config.TrainingMgmtChannel)

		trainingService.ISpreadsheetService.WriteCell(trainingService.config.TrainingSheet, rowNumber, trainingService.config.TrainingUtensilsColumn, "POSTED")
		glog.Infof("Updated event with trainings utensils=POSTED in sheet=%s", topic.sheet)

	}

	return nil
}

func (trainingService TrainingService) createTrainingMgmtPost(row []interface{}, params trainingParameters) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("Ban Training heint (")
	buffer.WriteString(row[trainingService.config.NameColumn].(string))
	buffer.WriteString(" - ")
	buffer.WriteString(row[trainingService.config.DateColumn].(string))
	buffer.WriteString(")")
	buffer.WriteString(" sein insgesomt *")
	buffer.WriteString(params.TotalGoing)
	buffer.WriteString("*, *")
	buffer.WriteString(params.GoingSGS07)
	buffer.WriteString(" SGS07* und *")
	buffer.WriteString(params.GoingSGS16)
	buffer.WriteString(" SGS16*.\n")

	if params.ResponsibleTrainingUtensils != "" {
		buffer.WriteString("Für Trainingsutensilien zuständig: *")
		buffer.WriteString(params.ResponsibleTrainingUtensils)
		buffer.WriteString("!*")
	}

	return buffer
}
