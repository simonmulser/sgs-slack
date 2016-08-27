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
  
  "google.golang.org/api/sheets/v4"
  "github.com/nlopes/slack"
)

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
    "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func toAlphabetChar(i int) string {
    return arr[i]
}

const TRAINING_SHEET = "1sKvrduUoKRfVk1QMaXwvps0IgAu6xUHhJ6XGSqJW68c"

func createSlackClient(slack_key string) *slack.Client{
  slackClient := slack.New(slack_key);
  // slackClient.SetDebug(true)

  return slackClient;
}

func createTrainingPost(row []interface{}) string {
  var buffer bytes.Buffer
  buffer.WriteString(row[0].(string))
  buffer.WriteString(" am *")
  buffer.WriteString(row[2].(string))
  buffer.WriteString(" um *")
  buffer.WriteString(row[3].(string))
  buffer.WriteString("*! Wer isch dobei?")
  return buffer.String()
 }

func writeCell(service *sheets.Service, row int, column int, text string) {
  var update_columns [][]interface{}
  var  update_rows []interface{}
  update_columns = append(update_columns, append(update_rows, text)) 

  valueRange := sheets.ValueRange{Values: update_columns}

  request := service.Spreadsheets.Values.Update(TRAINING_SHEET, 
    toAlphabetChar(column) + strconv.Itoa(row) + ":" + toAlphabetChar(column) + strconv.Itoa(row) , &valueRange)

  request.ValueInputOption("RAW")
  request.Do()
}

func main() {
  config := config.Read()
  service := google.New();

  readRange := "A2:F"
  response, error := service.Spreadsheets.Values.Get(TRAINING_SHEET, readRange).Do()
  if error != nil {
    log.Fatalf("Unable to retrieve data from sheet. %v", error)
  }

  // fmt.Print(response)
  fmt.Println(reflect.TypeOf(service))

  slackClient := createSlackClient(config.SLACK_KEY)

  if len(response.Values) > 0 {
    i := 2
    for _, row := range response.Values {

      if error != nil {
        log.Fatalf("Unable to update data from sheet. %v", error)
      }
      if(row[config.STATUS_COLUMN] == "POSTED"){
        postingDate, _ := time.Parse("02.01.2006", row[config.POSTING_DATE_COLUMN].(string))
        if(time.Now().After(postingDate)){
          fmt.Println(postingDate)
          params := slack.NewPostMessageParameters()
          params.AsUser = true
          message := createTrainingPost(row)
          channelId, timestamp, error := slackClient.PostMessage("test", message, params)

          if error != nil {
            log.Fatalf("Unable to update data from sheet. %v", error)
          }
          writeCell(service, i, config.STATUS_COLUMN, "POSTED")          
          writeCell(service, i, config.CHANNEL_ID_COLUMN, channelId)          
          writeCell(service, i, config.TIMESTAMP_COLUMN, timestamp)          
          }
        }

        i++
      }
    } else {
      fmt.Print("No data found.")
    }
}