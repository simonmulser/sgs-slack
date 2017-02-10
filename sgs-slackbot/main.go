package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
)

// Main holds all services
type Main struct {
	config             *Config
	slackService       *SlackService
	spreadsheetService *SpreadsheetService
	messageBuilder     *MessageBuilder
	gameService        *GameService
	trainingService    *TrainingService
}

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	glog.Info("Program started")

	env := "development"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	instance := Main{}
	instance.config = read(env)
	instance.slackService = newSlackService(instance.config.SlackKey)
	instance.messageBuilder = newMessageBuilder(instance.config, instance.slackService)
	instance.spreadsheetService = newSpreadsheetService()
	instance.trainingService = newTrainingService(&instance)
	instance.gameService = newGameService(&instance)

	instance.run()

	glog.Info("Program terminated")
}

func (main Main) run() {
	main.trainingService.process()

	main.gameService.process()
}
