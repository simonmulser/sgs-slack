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
	eventService       *EventService
	IMessageBuilder
	ISlackService
	ISpreadsheetService
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
	instance.ISlackService = newSlackService(instance.config.SlackKey)
	instance.IMessageBuilder = newMessageBuilder(instance.config, instance.slackService)
	instance.ISpreadsheetService = newSpreadsheetService()
	instance.eventService = newEventService(&instance)

	instance.run()

	glog.Info("Program terminated")
}

func (main Main) run() {
	main.eventService.process()
}
