package main

import (
  "bytes"
  "github.com/golang/glog"
  "time"
)

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