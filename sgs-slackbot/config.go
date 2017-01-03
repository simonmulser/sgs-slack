package main

import (
	"os"
	"encoding/json"
	"fmt"
)

type Config struct{
  TRAINING_SHEET string
  GAMES_07_SHEET string
  SLACK_KEY string

  TRAINING_CHANNEL string
  TRAINING_MGMT_CHANNEL string
  GAMES_07_CHANNEL string
  GAMES_16_CHANNEL string

  NAME_COLUMN int
  DESCRIPTION_COLUMN int
  DATE_COLUMN int
  POSTING_DATE_COLUMN int
  CHANNEL_ID_COLUMN int
  TIMESTAMP_COLUMN int
  TRAINING_UTENSILS_COLUMN int
  TRAINING_UTENSILS_RESPONSIBLE_TEXT string

  GAME_DATE_COLUMN int
  GAME_POSTING_DATE_COLUMN int
  HOME_COLUMN int
  AWAY_COLUMN int
  LOCATION_COLUMN int
  LOCATION_MAPS_COLUMN int
  SURFACE_COLUMN int
  GAME_CHANNEL_ID_COLUMN int
  GAME_TIMESTAMP_COLUMN int
}

func Read(env string) *Config{
	file, _ := os.Open("config/config.json")
	decoder := json.NewDecoder(file)
  config := Config{}
  err := decoder.Decode(&config)
  if err != nil {
    fmt.Println("error:", err)
  }

  file, _ = os.Open("config/connections/slack-key.json")
  decoder = json.NewDecoder(file)
  err = decoder.Decode(&config)
  if err != nil {
    fmt.Println("error:", err)
  }

  if (env == "development"){
    file, _ = os.Open("config/connections/development-config.json")
  } else if (env == "production"){
    file, _ = os.Open("config/connections/production-config.json")
  } else {
    fmt.Println("error: unkown env")
  }
  decoder = json.NewDecoder(file)
  err = decoder.Decode(&config)
  if err != nil {
    fmt.Println("error:", err)
  }

  return &config
}
