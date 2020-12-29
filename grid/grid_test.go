package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidInputsProduceErrors(t *testing.T) {

	inputs := []string{
		"",
		"A",
		"1",
		"A0",
		"A5",
		"E1",
		"Z1",
		"a1",
		"AA1",
	}

	grid := Grid{}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			_, _, err := grid.GetCoordinates(input)
			assert.Error(t, err)
		})
	}
}

type CoordinatesTestData struct {
	input       string
	rowIndex    int
	columnIndex int
}

func TestGetCoordinatesCovertsStringsToCoordinates(t *testing.T) {

	testData := []CoordinatesTestData{
		{"A1", 1, 1},
		{"A4", 4, 1},
		{"D1", 1, 4},
		{"D4", 4, 4},
	}

	grid := Grid{}

	for _, data := range testData {
		t.Run(data.input, func(t *testing.T) {
			rowIndex, columnIndex, err := grid.GetCoordinates(data.input)

			assert.NoError(t, err)
			assert.Equal(t, data.rowIndex, rowIndex)
			assert.Equal(t, data.columnIndex, columnIndex)
		})
	}
}
