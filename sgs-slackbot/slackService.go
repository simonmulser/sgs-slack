package main

import (
	"github.com/nlopes/slack"
)

type ISlackService interface {
	postMessage(channel string, message string) (string, string, error)
	updateMessage(channel, timestamp, text string) (string, string, string, error)
	getUserInfo(user string) (*slack.User, error)
}

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

func (slackService SlackService) updateMessage(channel, timestamp, text string) (string, string, string, error) {
	return slackService.slack.UpdateMessage(channel, timestamp, text)
}

func (slackService SlackService) getUserInfo(user string) (*slack.User, error) {
	return slackService.slack.GetUserInfo(user)
}
