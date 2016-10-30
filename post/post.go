package main

import (
  "github.com/golang/glog"
  "strconv"
  "time"
  "math/rand"
  "os"
  "flag"
  
  "github.com/nlopes/slack"
  "google.golang.org/api/sheets/v4"
)

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
    "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func toAlphabetChar(i int) string {
    return arr[i]
}

type TrainingParameters struct{
  Total_going string
  Going_sgs07 string
  Going_sgs16 string
  Responsible_balls string
}

type Main struct{
  config *Config
  slackClient *slack.Client
  service *sheets.Service
}

 func (main Main) createTrainingParams(reactions []slack.ItemReaction) TrainingParameters {
  var params TrainingParameters
  var going []string
  count_muscle := 0
  count_facepunch := 0

  for _,reaction := range reactions {
    if(reaction.Name == "muscle"){
      count_muscle= reaction.Count
      going = append(going, reaction.Users...)
      }
    if(reaction.Name == "facepunch"){
      count_facepunch = reaction.Count
      going = append(going, reaction.Users...)
      }
    }

  params.Going_sgs07 = strconv.Itoa(count_muscle)
  params.Going_sgs16 = strconv.Itoa(count_facepunch)
  params.Total_going = strconv.Itoa(count_muscle + count_facepunch)

  if(len(going) > 0){
    user, error := main.slackClient.GetUserInfo(going[rand.Intn(len(going))])
    if error != nil {
      glog.Fatalf("error: ", error.Error())
    }
    params.Responsible_balls = user.Name    
  }

  return params
 }

func (main Main) timeNow() time.Time {
  t := time.Now()
  utc, error := time.LoadLocation("Europe/Vienna")
  if error != nil {
    glog.Info("error: ", error.Error())
  }
  return t.In(utc)
}

func main() {
  flag.Parse()
  flag.Lookup("logtostderr").Value.Set("true")
  glog.Info("Program started")

  instance := Main{}
  instance.run()

  glog.Info("Program terminated")
}

func (main Main) run() {
  env := "development"
  if (len(os.Args) > 1) {
    env = os.Args[1]
  }
  main.config = Read(env)

  main.service = New();

  main.slackClient = main.createSlackClient(main.config.SLACK_KEY)

  var rows = main.readRange(main.config.TRAINING_SHEET ,"A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.CHANNEL_ID_COLUMN] == "FALSE" && main.timeNow().After(postingDate)){
        message := main.createTrainingPost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())
        if error != nil {
          glog.Fatalf("Unable to post message. %v", error)
        }

        main.writeCell(main.config.TRAINING_SHEET, i, main.config.CHANNEL_ID_COLUMN, channelId)
        main.writeCell(main.config.TRAINING_SHEET, i, main.config.TIMESTAMP_COLUMN, timestamp)
      }

      i++
      }
    } else {
      glog.Info("No data found.")
    }

  rows = main.readRange(main.config.TRAINING_SHEET, "A2:G")

  if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      date, error := time.Parse("02.01.2006 15:04", row[main.config.DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }
      date = date.Add(-8 * 60 * time.Minute)

      if(row[main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[main.config.BALLS_COLUMN] == "FALSE" && main.timeNow().After(date)){
          reactions, error := main.slackClient.GetReactions(
            slack.ItemRef{Channel: row[main.config.CHANNEL_ID_COLUMN].(string), Timestamp: row[main.config.TIMESTAMP_COLUMN].(string)},
            slack.GetReactionsParameters{})
          if error != nil {
            glog.Fatalf("Unable to get reactions. %v", error)
          }

          params := main.createTrainingParams(reactions)
          message := main.createTrainingMgmtPost(row, params)
          main.postMessage(main.config.TRAINING_MGMT_CHANNEL, message.String())
          main.postMessage("@" + params.Responsible_balls, main.config.BALLS_RESPONSIBLE_TEXT)
          main.writeCell(main.config.TRAINING_SHEET, i, main.config.BALLS_COLUMN, "TRUE")
          }
      i++
      }
    } else {
      glog.Info("No data found.")
    }

if false {
  rows = main.readRange(main.config.GAMES_07_SHEET, "A2:K")

    if len(rows.Values) > 0 {
    i := 2
    for _, row := range rows.Values {
      postingDate, error := time.Parse("02.01.2006 15:04", row[main.config.GAME_POSTING_DATE_COLUMN].(string))
      if error != nil {
        glog.Fatalf("Unable to parse date. %v", error)
      }

      if(row[main.config.GAME_CHANNEL_ID_COLUMN] == "FALSE" && main.timeNow().After(postingDate)){
        message := main.createGamePost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())
        if error != nil {
          glog.Fatalf("Unable to post massage. %v", error)
        }

        main.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_CHANNEL_ID_COLUMN, channelId)
        main.writeCell(main.config.GAMES_07_SHEET, i, main.config.GAME_TIMESTAMP_COLUMN, timestamp)
      }

      i++
      }
    } else {
      glog.Info("No data found.")
    }
  }
}
