package main

import (
  "github.com/golang/glog"
  "os"
  "flag"
)

type Main struct{
  config *Config
  slackService *SlackService
  spreadsheetService *SpreadsheetService
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
  instance.spreadsheetService = NewSpreadsheetService()

  instance.run()

  glog.Info("Program terminated")
}

func (main Main) run() {
  processTrainings(main)

  processGames(main)
}
