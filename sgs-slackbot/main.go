package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/simonmulser/config"
	"github.com/simonmulser/google"
	"github.com/simonmulser/slack"
)

// Main holds all services
type Main struct {
	config       *config.Config
	eventService *EventService
	IMessageBuilder
	slack.ISlackService
	google.ISpreadsheetService
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
	instance.config = config.Read(env)
	instance.ISlackService = slack.NewSlackService(instance.config.SlackKey)
	instance.ISpreadsheetService = google.NewSpreadsheetService()
	instance.eventService = newEventService(&instance)

	instance.run()

	glog.Info("Program terminated")
}

func (main Main) run() {
	main.eventService.process()
}
