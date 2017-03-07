package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/glog"
	nlopesslack "github.com/nlopes/slack"
	"github.com/simonmulser/config"
	"github.com/simonmulser/google"
	"github.com/simonmulser/slack"
)

type Message struct {
	Text string `json:"text"`
}

type Main struct {
	config *config.Config
	google.ISpreadsheetService
	slack.ISlackService
	topicConfig []TopicConfig
}

type TopicConfig struct {
	sheet string
	name  string
}

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	glog.Info("Program started")

	env := "development"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	instance := Main{}
	instance.config = config.Read(env)
	instance.ISpreadsheetService = google.NewSpreadsheetService()
	instance.ISlackService = slack.NewSlackService(instance.config.SlackKey)

	training := TopicConfig{instance.config.TrainingSheet, "Trainig"}
	sg07 := TopicConfig{instance.config.Games07Sheet, "Spiel 07"}
	sg16 := TopicConfig{instance.config.Games16Sheet, "Spiel 16"}

	instance.topicConfig = []TopicConfig{training, sg07, sg16}

	http.HandleFunc("/teilnehmer", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, instance.getSummary())
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}

func (instance Main) getSummary() string {
	var summary string

	for _, topic := range instance.topicConfig {
		readRange := instance.ISpreadsheetService.ReadRange(topic.sheet, "A2:L")

		firstPostedRow, error := instance.getFirstPostedRow(readRange.Values)

		if error != nil {
			glog.Warningf("Could not find first posted row in sheet=%s", topic.sheet)
			continue
		}
		reactions, error := instance.ISlackService.GetReactions(
			nlopesslack.ItemRef{Channel: firstPostedRow[instance.config.ChannelIDColumn].(string), Timestamp: firstPostedRow[instance.config.TimestampColumn].(string)},
			nlopesslack.GetReactionsParameters{})

		if error != nil {
			glog.Warningf("Could not retrieve reactions for row in sheet=%s", topic.sheet)
			continue
		}

		summary += instance.createEventReactionsSummary(topic, firstPostedRow, reactions)
	}

	return summary
}

func (main Main) getFirstPostedRow(rows [][]interface{}) ([]interface{}, error) {
	for i, row := range rows {
		if row[main.config.StatusColumn] == "POSTED" {
			glog.Infof("Found POSTED row in row=%v", i)
			return row, nil
		}
	}
	return nil, errors.New("No row with status 'POSTED' found")
}

func (main Main) createEventReactionsSummary(topic TopicConfig, row []interface{}, reactions []nlopesslack.ItemReaction) string {
	var summary string

	summary += "*" + topic.name + " " + row[main.config.DateColumn].(string) + "* "
	if len(reactions) == 0 {
		summary += "keine Reaktionen bis jetzt."
	} else {
		for i, reaction := range reactions {
			summary += ":" + reaction.Name + ": " + strconv.Itoa(reaction.Count)
			if i != (len(reactions) - 1) {
				summary += ", "
			} else {
				summary += "\n"
			}
		}
	}

	return summary
}
