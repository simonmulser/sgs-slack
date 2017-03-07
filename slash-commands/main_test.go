package main

import (
	"os"
	"testing"

	"github.com/simonmulser/config"
	"github.com/stretchr/testify/assert"
)

var mainInstance Main

func TestMain(m *testing.M) {
	mainInstance = Main{}
	mainInstance.config = createConfig()

	os.Exit(m.Run())
}

func TestGetFirstPostedRow1(t *testing.T) {
	firstRow := createRow([]string{"OVER", "first"})
	secondRow := createRow([]string{"POSTED", "second"})
	thirdRow := createRow([]string{"POSTED", "third"})

	var rows [][]interface{}
	rows = append(rows, firstRow)
	rows = append(rows, secondRow)
	rows = append(rows, thirdRow)
	row, error := mainInstance.getFirstPostedRow(rows)

	assert.NotNil(t, row)
	assert.Nil(t, error)
	assert.Equal(t, "second", row[1].(string))
}

func TestGetFirstPostedRow2(t *testing.T) {
	firstRow := createRow([]string{"POSTED", "first"})
	secondRow := createRow([]string{"NEW", "second"})
	thirdRow := createRow([]string{"OVER", "third"})

	var rows [][]interface{}
	rows = append(rows, firstRow)
	rows = append(rows, secondRow)
	rows = append(rows, thirdRow)
	row, error := mainInstance.getFirstPostedRow(rows)

	assert.NotNil(t, row)
	assert.Nil(t, error)
	assert.Equal(t, "first", row[1].(string))
}

func TestGetFirstPostedRowNowRows(t *testing.T) {
	var rows [][]interface{}
	row, error := mainInstance.getFirstPostedRow(rows)

	assert.Nil(t, row)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "with status 'POSTED'")
}

func TestGetFirstPostedRowNowRowsMatching(t *testing.T) {
	firstRow := createRow([]string{"NEW", "first"})
	secondRow := createRow([]string{"OVER", "second"})

	var rows [][]interface{}
	rows = append(rows, firstRow)
	rows = append(rows, secondRow)
	row, error := mainInstance.getFirstPostedRow(rows)

	assert.Nil(t, row)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "with status 'POSTED'")
}

func createRow(data []string) []interface{} {
	row := make([]interface{}, len(data))
	for i, s := range data {
		row[i] = s
	}

	return row
}

func createConfig() *config.Config {
	config := config.Config{
		StatusColumn: 0,
	}
	return &config
}
