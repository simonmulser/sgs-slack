package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/nlopes/slack"
)

type ITopicCommand interface {
	execute(row []interface{}, topic topicConfig, rowNumber int) error
}

// TrainingService processes the trainings
type TrainingService struct {
	config *Config
	ISlackService
	ISpreadsheetService
}

func newTrainingService(config *Config, slackService ISlackService, spreadsheetService ISpreadsheetService) *TrainingService {
	trainingService := new(TrainingService)
	trainingService.config = config
	trainingService.ISlackService = slackService
	trainingService.ISpreadsheetService = spreadsheetService

	return trainingService
}

func (trainingService TrainingService) execute(row []interface{}, topic topicConfig, rowNumber int) error {
	date, error := time.Parse("02.01.2006 15:04", row[trainingService.config.DateColumn].(string))
	if error != nil {
		glog.Warningf("Unable to parse date. %v", error)
		return error
	}
	date = date.Add(-8 * 60 * time.Minute)

	if row[trainingService.config.StatusColumn] == "POSTED" && row[trainingService.config.TrainingUtensilsColumn] != "POSTED" && timeNow().After(date) {
		reactions, error := trainingService.ISlackService.getReactions(
			slack.ItemRef{Channel: row[trainingService.config.ChannelIDColumn].(string), Timestamp: row[trainingService.config.TimestampColumn].(string)},
			slack.GetReactionsParameters{})

		if error != nil {
			glog.Warningf("Unable to get reactions. %v", error)
			return error
		}

		params := trainingService.createTrainingParams(reactions)
		message := trainingService.createTrainingMgmtPost(row, params)

		trainingService.ISlackService.postMessage(trainingService.config.TrainingMgmtChannel, message.String())
		trainingService.ISlackService.postMessage("@"+params.ResponsibleTrainingUtensils, trainingService.config.TrainingUtensilsResponsibleText)
		trainingService.ISpreadsheetService.writeCell(trainingService.config.TrainingSheet, rowNumber, trainingService.config.TrainingUtensilsColumn, "POSTED")
		glog.Info("selected responsible person and updated sheets")
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

func (trainingService TrainingService) createTrainingParams(reactions []slack.ItemReaction) trainingParameters {
	var params trainingParameters
	var going []string
	countMuscle := 0
	countFacepunch := 0

	for _, reaction := range reactions {
		if reaction.Name == "muscle" {
			countMuscle = reaction.Count
			going = append(going, reaction.Users...)
		}
		if reaction.Name == "facepunch" {
			countFacepunch = reaction.Count
			going = append(going, reaction.Users...)
		}
	}

	params.GoingSGS07 = strconv.Itoa(countMuscle)
	params.GoingSGS16 = strconv.Itoa(countFacepunch)
	params.TotalGoing = strconv.Itoa(countMuscle + countFacepunch)

	if len(going) > 0 {
		user, error := trainingService.ISlackService.getUserInfo(going[rand.Intn(len(going))])
		if error != nil {
			glog.Fatalf("error: %v", error)
		}
		params.ResponsibleTrainingUtensils = user.Name
	}

	return params
}
