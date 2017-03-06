package slack

import (
	"github.com/nlopes/slack"
)

type ISlackService interface {
	PostMessage(channel string, message string) (string, string, error)
	UpdateMessage(channel, timestamp, text string) (string, string, string, error)
	GetUserInfo(user string) (*slack.User, error)
	GetReactions(item slack.ItemRef, params slack.GetReactionsParameters) ([]slack.ItemReaction, error)
}

// SlackService helps you to send messages to slack
type SlackService struct {
	slack *slack.Client
}

func NewSlackService(SlackKey string) *SlackService {
	slackService := new(SlackService)
	slackService.slack = slack.New(SlackKey)
	// slackService.slack.SetDebug(true)

	return slackService
}

func (slackService SlackService) PostMessage(channel string, message string) (string, string, error) {
	params := slack.NewPostMessageParameters()
	params.AsUser = true
	return slackService.slack.PostMessage(channel, message, params)
}

func (slackService SlackService) UpdateMessage(channel, timestamp, text string) (string, string, string, error) {
	return slackService.slack.UpdateMessage(channel, timestamp, text)
}

func (slackService SlackService) GetUserInfo(user string) (*slack.User, error) {
	return slackService.slack.GetUserInfo(user)
}

func (slackService SlackService) GetReactions(item slack.ItemRef, params slack.GetReactionsParameters) ([]slack.ItemReaction, error) {
	return slackService.slack.GetReactions(item, params)
}
