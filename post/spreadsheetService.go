package main

import (
  "google.golang.org/api/sheets/v4"
  "github.com/golang/glog"
  "strconv"
)

type SpreadsheetService struct{
    service *sheets.Service
}

func NewSpreadsheetService() *SpreadsheetService{
  spreadsheetService := new(SpreadsheetService)
  spreadsheetService.service = New()

  return spreadsheetService
}

func (spreadsheetService SpreadsheetService) writeCell(sheet string, row int, column int, text string) {
  var update_columns [][]interface{}
  var  update_rows []interface{}
  update_columns = append(update_columns, append(update_rows, text)) 

  valueRange := sheets.ValueRange{Values: update_columns}

  request := spreadsheetService.service.Spreadsheets.Values.Update(sheet,
    toAlphabetChar(column) + strconv.Itoa(row) + ":" + toAlphabetChar(column) + strconv.Itoa(row) , &valueRange)

  request.ValueInputOption("RAW")
  request.Do()
}

func (spreadsheetService SpreadsheetService) readRange(sheet string, rangeToRead string) *sheets.ValueRange {
    response, error := spreadsheetService.service.Spreadsheets.Values.Get(sheet, rangeToRead).Do()
    if error != nil {
      glog.Fatalf("Unable to retrieve data from sheet. %v", error)
    }

    return response
}
