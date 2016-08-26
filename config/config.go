package config

import (
	"os"
	"encoding/json"
	"fmt"
)

type Config struct{
  TRAINING_SHEET string
  SLACK_KEY string
}

func Read() *Config{
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
  	config := Config{}
  	err := decoder.Decode(&config)
  	if err != nil {
    	fmt.Println("error:", err)
  	}
  	return &config
}