package main

import (
	"os"
	"testing"

	"github.com/simonmulser/testutils"
	"github.com/stretchr/testify/assert"
)

var mainInstance Main

func TestMain(m *testing.M) {
	mainInstance = Main{}
	mainInstance.config = testutils.CreateConfig()

	os.Exit(m.Run())
}

func TestGetFirstPostedRow1(t *testing.T) {
	firstRow := testutils.CreateRow([]string{"OVER", "first"})
	secondRow := testutils.CreateRow([]string{"POSTED", "second"})
	thirdRow := testutils.CreateRow([]string{"POSTED", "third"})

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
	firstRow := testutils.CreateRow([]string{"POSTED", "first"})
	secondRow := testutils.CreateRow([]string{"NEW", "second"})
	thirdRow := testutils.CreateRow([]string{"OVER", "third"})

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
	firstRow := testutils.CreateRow([]string{"NEW", "first"})
	secondRow := testutils.CreateRow([]string{"OVER", "second"})

	var rows [][]interface{}
	rows = append(rows, firstRow)
	rows = append(rows, secondRow)
	row, error := mainInstance.getFirstPostedRow(rows)

	assert.Nil(t, row)
	assert.NotNil(t, error)
	assert.Contains(t, error.Error(), "with status 'POSTED'")
}


}
