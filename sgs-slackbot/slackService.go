package main

import (
  "github.com/nlopes/slack"
)

type SlackService struct{
	slack *slack.Client
}

func NewSlackService(slack_key string) *SlackService{
  slackService := new(SlackService)
  slackService.slack = slack.New(slack_key);
  // slackService.slack.SetDebug(true)

  return slackService;
}

func (slackService SlackService) postMessage(channel string, message string) (string, string, error) {
  params := slack.NewPostMessageParameters()
  params.AsUser = true
  return slackService.slack.PostMessage(channel, message, params)
}