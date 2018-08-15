package config

import (
	"encoding/json"
	"os"

	"github.com/golang/glog"
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
	DayOfWeekColumn   int

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

const configFile = "./config/config.json" 
const slackKeyFile = "./config/connections/slack-key.json"
const developmentFile = "./config/connections/development-config.json"
const productionFile = "./config/connections/production-config.json"

func Read(env string) *Config {
	file, err := os.Open(configFile)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}

	file, err = os.Open(slackKeyFile)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}
	decoder = json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}

	fileName := ""
	switch env {
	case "development":
		fileName = developmentFile
	case "production":
		fileName = productionFile
	default:
		glog.Fatalf("unkown environment %s", env)
	}

	file, err = os.Open(fileName)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}
	decoder = json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		glog.Fatalf("error: %s", err)
	}

	return &config
}
