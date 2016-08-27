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

func writeCell(service *sheets.Service, row int, column string, text string) {
  var update_columns [][]interface{}
  var  update_rows []interface{}
  update_columns = append(update_columns, append(update_rows, text)) 

  valueRange := sheets.ValueRange{Values: update_columns}

  request := service.Spreadsheets.Values.Update(TRAINING_SHEET, 
    column + strconv.Itoa(row) + ":" + column + strconv.Itoa(row) , &valueRange)

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
      if(row[4] == "FALSE"){
        postingDate, _ := time.Parse("02.01.2006", row[5].(string))
        if(time.Now().After(postingDate)){
          fmt.Println(postingDate)
          params := slack.NewPostMessageParameters()
          params.AsUser = true
          message := createTrainingPost(row)
          slackClient.PostMessage("test", message, params)

          writeCell(service, i, "E","FALSE")          
          }
        }

        i++
      }
    } else {
      fmt.Print("No data found.")
    }
}