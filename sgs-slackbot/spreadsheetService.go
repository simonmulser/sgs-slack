package main

import (
	"strconv"

	"github.com/golang/glog"
	"google.golang.org/api/sheets/v4"
)

type ISpreadsheetService interface {
	writeCell(sheet string, row int, column int, text string)
	readRange(sheet string, rangeToRead string) *sheets.ValueRange
}

// SpreadsheetService helps you to interact with a sheet
type SpreadsheetService struct {
	service *sheets.Service
}

func newSpreadsheetService() *SpreadsheetService {
	spreadsheetService := new(SpreadsheetService)
	spreadsheetService.service = New()

	return spreadsheetService
}

func (spreadsheetService SpreadsheetService) writeCell(sheet string, row int, column int, text string) {
	var upDateColumns [][]interface{}
	var updateRows []interface{}
	upDateColumns = append(upDateColumns, append(updateRows, text))

	valueRange := sheets.ValueRange{Values: upDateColumns}

	request := spreadsheetService.service.Spreadsheets.Values.Update(sheet,
		toAlphabetChar(column)+strconv.Itoa(row)+":"+toAlphabetChar(column)+strconv.Itoa(row), &valueRange)

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
