package testutils

import "github.com/simonmulser/config"

func CreateRow(data []string) []interface{} {
	row := make([]interface{}, len(data))
	for i, s := range data {
		row[i] = s
	}

	return row
}

func CreateConfig() *config.Config {
	config := config.Config{
		StatusColumn:                    0,
		PostingDateColumn:               1,
		ChannelIDColumn:                 2,
		TimestampColumn:                 3,
		DateColumn:                      4,
		TrainingUtensilsColumn:          5,
		TrainingUtensilsResponsibleText: "FancyText",
		TrainingMgmtChannel:             "channelID",
	}
	return &config
}
