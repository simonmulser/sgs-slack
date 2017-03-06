package main

import (
	"math/rand"
	"strconv"

	"github.com/golang/glog"
	nlopesslack "github.com/nlopes/slack"
	"github.com/simonmulser/slack"
)

type ITrainingParamsService interface {
	create(reactions []nlopesslack.ItemReaction) trainingParameters
}

type TrainingParamsService struct {
	slack.ISlackService
}

func newTrainingParamsService(slackService slack.ISlackService) *TrainingParamsService {
	trainingService := new(TrainingParamsService)
	trainingService.ISlackService = slackService

	return trainingService
}

type trainingParameters struct {
	TotalGoing                  string
	GoingSGS07                  string
	GoingSGS16                  string
	ResponsibleTrainingUtensils string
}

func (trainingParamsService TrainingParamsService) create(reactions []nlopesslack.ItemReaction) trainingParameters {
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
		user, error := trainingParamsService.ISlackService.GetUserInfo(going[rand.Intn(len(going))])
		if error != nil {
			glog.Fatalf("error: %v", error)
		}
		params.ResponsibleTrainingUtensils = user.Name
	}

	return params
}
