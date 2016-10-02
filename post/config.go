package main

import (
	"os"
	"encoding/json"
	"fmt"
)

type Config struct{
  TRAINING_SHEET string
  SLACK_KEY string
  TRAINING_CHANNEL string
  TRAINING_MGMT_CHANNEL string
  NAME_COLUMN int
  DESCRIPTION_COLUMN int
  DATE_COLUMN int
  POSTING_DATE_COLUMN int
  CHANNEL_ID_COLUMN int
  TIMESTAMP_COLUMN int
  BALLS_COLUMN int
  BALLS_RESPONSIBLE_TEXT string
}

func Read(env string) *Config{
	file, _ := os.Open("config/config.json")
	decoder := json.NewDecoder(file)
  config := Config{}
  err := decoder.Decode(&config)
  if err != nil {
    fmt.Println("error:", err)
  }

  if (env == "development"){
    file, _ = os.Open("config/development-config.json")
  } else if (env == "production"){
    file, _ = os.Open("config/production-config.json")
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
