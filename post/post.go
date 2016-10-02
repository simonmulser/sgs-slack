package main

import (
  "fmt"
  "log"
  "bytes"
  "reflect"
  "strconv"
  "time"
  "math/rand"
  "os"
  
  "google.golang.org/api/sheets/v4"
  "github.com/nlopes/slack"
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
}

func (main Main) createSlackClient(slack_key string) *slack.Client{
  slackClient := slack.New(slack_key);
  // slackClient.SetDebug(true)

  return slackClient;
}

func (main Main) createTrainingPost(row []interface{}) bytes.Buffer {
  var buffer bytes.Buffer
  buffer.WriteString(row[main.config.NAME_COLUMN].(string))
  buffer.WriteString(" am *")
  buffer.WriteString(row[main.config.DATE_COLUMN].(string))
  buffer.WriteString("* Uhr! ")
  buffer.WriteString(row[main.config.DESCRIPTION_COLUMN].(string))
  return buffer
 }

 func (main Main) createGamePost(row []interface{}) bytes.Buffer {
  var buffer bytes.Buffer

  meetingTime, _ := time.Parse("02.01.2006 15:04", row[main.config.GAME_DATE_COLUMN].(string)) 
  meetingTime = meetingTime.Add(-1 * 60 * time.Minute)

  buffer.WriteString("*")
  buffer.WriteString(row[main.config.HOME_COLUMN].(string))
  buffer.WriteString(" : ")
  buffer.WriteString(row[main.config.AWAY_COLUMN].(string))
  buffer.WriteString("* am *")
  buffer.WriteString(row[main.config.GAME_DATE_COLUMN].(string))
  buffer.WriteString("* Uhr auf ")
  buffer.WriteString(row[main.config.SURFACE_COLUMN].(string))
  buffer.WriteString("! Treffpunkt: ")
  buffer.WriteString(meetingTime.Format("15:04"))
  buffer.WriteString(", ")
  buffer.WriteString(row[main.config.LOCATION_COLUMN].(string))
  buffer.WriteString(" - ")
  buffer.WriteString(row[main.config.LOCATION_MAPS_COLUMN].(string))
  buffer.WriteString(".")

  return buffer
 }


func (main Main) createTrainingMgmtPost(row []interface{}, params TrainingParameters ) bytes.Buffer {
  var buffer bytes.Buffer
  buffer.WriteString("Ban Training heint (")
  buffer.WriteString(row[main.config.NAME_COLUMN].(string))
  buffer.WriteString(" - ")
  buffer.WriteString(row[main.config.DATE_COLUMN].(string))
  buffer.WriteString(")")
  buffer.WriteString(" sein insgesomt *")
  buffer.WriteString(params.Total_going)
  buffer.WriteString("*, *");
  buffer.WriteString(params.Going_sgs07)
  buffer.WriteString(" SGS07* und *");
  buffer.WriteString(params.Going_sgs16)
  buffer.WriteString(" SGS16*.\n");

  if(params.Responsible_balls != "") {
    buffer.WriteString("Für die Bälle zuständig: *");
    buffer.WriteString(params.Responsible_balls)
    buffer.WriteString("!*");
  }

  return buffer
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
    user, _ := main.slackClient.GetUserInfo(going[rand.Intn(len(going))])
    params.Responsible_balls = user.Name    
  }

  return params
 }

func (main Main) writeCell(service *sheets.Service, sheet string, row int, column int, text string) {
  var update_columns [][]interface{}
  var  update_rows []interface{}
  update_columns = append(update_columns, append(update_rows, text)) 

  valueRange := sheets.ValueRange{Values: update_columns}

  request := service.Spreadsheets.Values.Update(sheet, 
    toAlphabetChar(column) + strconv.Itoa(row) + ":" + toAlphabetChar(column) + strconv.Itoa(row) , &valueRange)

  request.ValueInputOption("RAW")
  request.Do()
}

func (main Main) postMessage(channel string, message string) (string, string, error) {
  params := slack.NewPostMessageParameters()
  params.AsUser = true
  return main.slackClient.PostMessage(channel, message, params)
}

func (main Main) timeNow() time.Time {
  t := time.Now()
  utc, err := time.LoadLocation("Europe/Vienna")
  if err != nil {
    fmt.Println("err: ", err.Error())
  }
  return t.In(utc)
}

func main() {
  instance := Main{}
  instance.run()
}

func (main Main) run() {
  env := "development"
  if (len(os.Args) > 1) {
    env = os.Args[1]
  }
  main.config = Read(env)

  service := New();
  readRange := "A2:J"

  fmt.Println(reflect.TypeOf(service))

  main.slackClient = main.createSlackClient(main.config.SLACK_KEY)

  response, error := service.Spreadsheets.Values.Get(main.config.TRAINING_SHEET, readRange).Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }
  if len(response.Values) > 0 {
    i := 2
    for _, row := range response.Values {
      postingDate, _ := time.Parse("02.01.2006 15:04", row[main.config.POSTING_DATE_COLUMN].(string))

      if(row[main.config.CHANNEL_ID_COLUMN] == "FALSE" && main.timeNow().After(postingDate)){
        message := main.createTrainingPost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())

        if error != nil {
          log.Fatalf("Unable to update data from sheet. %v", error)
        }
        main.writeCell(service, main.config.TRAINING_SHEET, i, main.config.CHANNEL_ID_COLUMN, channelId)          
        main.writeCell(service, main.config.TRAINING_SHEET, i, main.config.TIMESTAMP_COLUMN, timestamp)          
      }

      i++
      }
    } else {
      fmt.Print("No data found.")
    }

  response, error = service.Spreadsheets.Values.Get(main.config.TRAINING_SHEET, readRange).Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }
  if len(response.Values) > 0 {
    i := 2
    for _, row := range response.Values {
      date, _ := time.Parse("02.01.2006 15:04", row[main.config.DATE_COLUMN].(string))
      date = date.Add(-8 * 60 * time.Minute)

      if(row[main.config.CHANNEL_ID_COLUMN] != "FALSE" && row[main.config.BALLS_COLUMN] == "FALSE" && main.timeNow().After(date)){
          reactions, error := main.slackClient.GetReactions(
            slack.ItemRef{Channel: row[main.config.CHANNEL_ID_COLUMN].(string), Timestamp: row[main.config.TIMESTAMP_COLUMN].(string)},
            slack.GetReactionsParameters{})
          if error != nil {
            log.Fatalf("Unable to update data from sheet. %v", error)
          }

          params := main.createTrainingParams(reactions)
          message := main.createTrainingMgmtPost(row, params)
          main.postMessage(main.config.TRAINING_MGMT_CHANNEL, message.String())
          main.postMessage("@" + params.Responsible_balls, main.config.BALLS_RESPONSIBLE_TEXT)
          main.writeCell(service, main.config.TRAINING_SHEET, i, main.config.BALLS_COLUMN, "TRUE")          
          }
      i++
      }
    } else {
      fmt.Print("No data found.")
    }

  response, error = service.Spreadsheets.Values.Get(main.config.GAMES_07_SHEET, "A2:K").Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }
  if false {
    i := 2
    for _, row := range response.Values {
      postingDate, _ := time.Parse("02.01.2006 15:04", row[main.config.GAME_POSTING_DATE_COLUMN].(string))

      if(row[main.config.GAME_CHANNEL_ID_COLUMN] == "FALSE" && main.timeNow().After(postingDate)){
        message := main.createGamePost(row)
        channelId, timestamp, error := main.postMessage(main.config.TRAINING_CHANNEL, message.String())

        if error != nil {
          log.Fatalf("Unable to update data from sheet. %v", error)
        }
        main.writeCell(service, main.config.GAMES_07_SHEET, i, main.config.GAME_CHANNEL_ID_COLUMN, channelId)
        main.writeCell(service, main.config.GAMES_07_SHEET, i, main.config.GAME_TIMESTAMP_COLUMN, timestamp)
      }

      i++
      }
    } else {
      fmt.Print("No data found.")
    }
}
