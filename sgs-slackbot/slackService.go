package main

import (
	"github.com/nlopes/slack"
)

// SlackService helps you to send messages to slack
type SlackService struct {
	slack *slack.Client
}

func newSlackService(SlackKey string) *SlackService {
	slackService := new(SlackService)
	slackService.slack = slack.New(SlackKey)
	// slackService.slack.SetDebug(true)

	return slackService
}

func (slackService SlackService) postMessage(channel string, message string) (string, string, error) {
	params := slack.NewPostMessageParameters()
	params.AsUser = true
	return slackService.slack.PostMessage(channel, message, params)
}
