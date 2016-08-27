package main

import (
  "github.com/simonmulser/config"
  "github.com/simonmulser/google"

  "fmt"
  "log"
  "bytes"
  "reflect"
  "strconv"
  "time"
  "math/rand"
  
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

func createSlackClient(slack_key string) *slack.Client{
  slackClient := slack.New(slack_key);
  // slackClient.SetDebug(true)

  return slackClient;
}

func createTrainingPost(row []interface{}) bytes.Buffer {
  var buffer bytes.Buffer
  buffer.WriteString(row[0].(string))
  buffer.WriteString(" am *")
  buffer.WriteString(row[2].(string))
  buffer.WriteString(" um ")
  buffer.WriteString(row[3].(string))
  buffer.WriteString("*! Wer isch dobei?")
  return buffer
 }

func updateTrainingPost(row []interface{}, params TrainingParameters ) bytes.Buffer {
  buffer := createTrainingPost(row)

  buffer.WriteString("\nEs sein insgesomt *")
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

 func createTrainingParams(reactions []slack.ItemReaction, slackClient *slack.Client) TrainingParameters {
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
    user, _ := slackClient.GetUserInfo(going[rand.Intn(len(going))])
    params.Responsible_balls = user.Name    
  }

  return params
 }

func writeCell(service *sheets.Service, sheet string, row int, column int, text string) {
  var update_columns [][]interface{}
  var  update_rows []interface{}
  update_columns = append(update_columns, append(update_rows, text)) 

  valueRange := sheets.ValueRange{Values: update_columns}

  request := service.Spreadsheets.Values.Update(sheet, 
    toAlphabetChar(column) + strconv.Itoa(row) + ":" + toAlphabetChar(column) + strconv.Itoa(row) , &valueRange)

  request.ValueInputOption("RAW")
  request.Do()
}

func main() {
  config := config.Read()
  service := google.New();
  readRange := "A2:J"

  fmt.Println(reflect.TypeOf(service))

  slackClient := createSlackClient(config.SLACK_KEY)

  response, error := service.Spreadsheets.Values.Get(config.TRAINING_SHEET, readRange).Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }
  if len(response.Values) > 0 {
    i := 2
    for _, row := range response.Values {
      postingDate, _ := time.Parse("02.01.2006", row[config.POSTING_DATE_COLUMN].(string))

      if(row[config.STATUS_COLUMN] == "FALSE" && time.Now().After(postingDate)){
        fmt.Println(postingDate)
        params := slack.NewPostMessageParameters()
        params.AsUser = true
        message := createTrainingPost(row)
        channelId, timestamp, error := slackClient.PostMessage("test", message.String(), params)

        if error != nil {
          log.Fatalf("Unable to update data from sheet. %v", error)
        }
        writeCell(service, config.TRAINING_SHEET, i, config.STATUS_COLUMN, "TRUE")          
        writeCell(service, config.TRAINING_SHEET, i, config.CHANNEL_ID_COLUMN, channelId)          
        writeCell(service, config.TRAINING_SHEET, i, config.TIMESTAMP_COLUMN, timestamp)          
        writeCell(service, config.TRAINING_SHEET, i, config.BALLS_COLUMN, "FALSE")          
      }

      i++
      }
    } else {
      fmt.Print("No data found.")
    }

  response, error = service.Spreadsheets.Values.Get(config.TRAINING_SHEET, readRange).Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }
  if len(response.Values) > 0 {
    i := 2
    for _, row := range response.Values {
      date, _ := time.Parse("02.01.2006", row[config.DATE_COLUMN].(string))
      
      if(row[config.STATUS_COLUMN] == "TRUE" && row[config.BALLS_COLUMN] == "FALSE" && time.Now().After(date)){
          reactions, error := slackClient.GetReactions(
            slack.ItemRef{Channel: row[config.CHANNEL_ID_COLUMN].(string), Timestamp: row[config.TIMESTAMP_COLUMN].(string)},
            slack.GetReactionsParameters{})
          if error != nil {
            log.Fatalf("Unable to update data from sheet. %v", error)
          }

          params := createTrainingParams(reactions, slackClient)
          message := updateTrainingPost(row, params)
          slackClient.UpdateMessage(row[config.CHANNEL_ID_COLUMN].(string),
            row[config.TIMESTAMP_COLUMN].(string), message.String())
          }
      i++
      }
    } else {
      fmt.Print("No data found.")
    }
}