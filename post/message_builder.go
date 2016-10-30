package main

import (
  "bytes"
  "strconv"
  "time"
  "math/rand"

  "github.com/golang/glog"
  "github.com/nlopes/slack"
)

type TrainingParameters struct{
  Total_going string
  Going_sgs07 string
  Going_sgs16 string
  Responsible_balls string
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

  meetingTime, error := time.Parse("02.01.2006 15:04", row[main.config.GAME_DATE_COLUMN].(string)) 
  if error != nil {
    glog.Fatalf("Unable to parse date. %v", error)
  }
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
    user, error := main.slackClient.GetUserInfo(going[rand.Intn(len(going))])
    if error != nil {
      glog.Fatalf("error: ", error.Error())
    }
    params.Responsible_balls = user.Name
  }

  return params
 }