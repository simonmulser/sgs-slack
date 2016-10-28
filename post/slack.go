package main

import (
  "github.com/nlopes/slack"
)

func (main Main) createSlackClient(slack_key string) *slack.Client{
  slackClient := slack.New(slack_key);
  // slackClient.SetDebug(true)

  return slackClient;
}

func (main Main) postMessage(channel string, message string) (string, string, error) {
  params := slack.NewPostMessageParameters()
  params.AsUser = true
  return main.slackClient.PostMessage(channel, message, params)
}