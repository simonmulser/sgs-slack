package main

import (
  "github.com/nlopes/slack"
)

type SlackClient struct{
	slack *slack.Client
}

func NewSlackClient(slack_key string) *SlackClient{
  slackClient := new(SlackClient)
  slackClient.slack = slack.New(slack_key);
  // slackClient.slack.SetDebug(true)

  return slackClient;
}

func (slackClient SlackClient) postMessage(channel string, message string) (string, string, error) {
  params := slack.NewPostMessageParameters()
  params.AsUser = true
  return slackClient.slack.PostMessage(channel, message, params)
}