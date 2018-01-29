package google

import (
	"strconv"

	"github.com/golang/glog"
	"google.golang.org/api/sheets/v4"
)

type ISpreadsheetService interface {
	WriteAll(sheet string, row int, status string, channel string, timestamp string)
	WriteCell(sheet string, row int, column int, text string)
	ReadRange(sheet string, rangeToRead string) *sheets.ValueRange
}

// SpreadsheetService helps you to interact with a sheet
type SpreadsheetService struct {
	service *sheets.Service
}

func NewSpreadsheetService() *SpreadsheetService {
	spreadsheetService := new(SpreadsheetService)
	spreadsheetService.service = New()

	return spreadsheetService
}

func (spreadsheetService SpreadsheetService) WriteAll(sheet string, row int, status string, channel string, timestamp string) {
	var upDateColumns [][]interface{}
	var updateRows []interface{}
	upDateColumns = append(upDateColumns, append(updateRows, status, channel, timestamp))

	valueRange := sheets.ValueRange{Values: upDateColumns}

	request := spreadsheetService.service.Spreadsheets.Values.Update(sheet,
		"A"+strconv.Itoa(row)+":"+"C"+strconv.Itoa(row), &valueRange)

	request.ValueInputOption("RAW")
	_, error := request.Do()
	if error != nil {
		glog.Fatalf("Unable to update data. %v", error)
	}
}

func (spreadsheetService SpreadsheetService) WriteCell(sheet string, row int, column int, text string) {
	var upDateColumns [][]interface{}
	var updateRows []interface{}
	upDateColumns = append(upDateColumns, append(updateRows, text))

	valueRange := sheets.ValueRange{Values: upDateColumns}

	request := spreadsheetService.service.Spreadsheets.Values.Update(sheet,
		toAlphabetChar(column)+strconv.Itoa(row)+":"+toAlphabetChar(column)+strconv.Itoa(row), &valueRange)

	request.ValueInputOption("RAW")
	_, error := request.Do()
	if error != nil {
		glog.Fatalf("Unable to update data. %v", error)
	}
}

func (spreadsheetService SpreadsheetService) ReadRange(sheet string, rangeToRead string) *sheets.ValueRange {
	response, error := spreadsheetService.service.Spreadsheets.Values.Get(sheet, rangeToRead).Do()
	if error != nil {
		glog.Fatalf("Unable to retrieve data from sheet. %v", error)
	}

	return response
}

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func toAlphabetChar(i int) string {
	return arr[i]
}
