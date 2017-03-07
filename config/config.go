package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config contains all config data
type Config struct {
	TrainingSheet string
	Games07Sheet  string
	Games16Sheet  string
	SlackKey      string

	TrainingChannel     string
	TrainingMgmtChannel string
	Games07Channel      string
	Games16Channel      string

	StatusColumn      int
	PostingDateColumn int
	ChannelIDColumn   int
	TimestampColumn   int
	DateColumn        int

	TrainingUtensilsColumn          int
	NameColumn                      int
	DescriptionColumn               int
	TrainingUtensilsResponsibleText string

	HomeColumn         int
	AwayColumn         int
	LocationColumn     int
	LocationMapsColumn int
	SurfaceColumn      int
}

func Read(env string) *Config {
	file, _ := os.Open("../config/config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	file, _ = os.Open("../config/connections/slack-key.json")
	decoder = json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	if env == "development" {
		file, _ = os.Open("../config/connections/development-config.json")
	} else if env == "production" {
		file, _ = os.Open("../config/connections/production-config.json")
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