package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CoordinatesTestData struct {
	input           string
	rowIndex        int
	columnIndex     int
	shouldBeAnError bool
}

func TestGetCoordinates(t *testing.T) {

	testData := []CoordinatesTestData{
		{"", 0, 0, true},
	}

	grid := grid.Grid{}

	for _, data := range testData {
		t.Run(data.input, func(t *testing.T) {
			rowIndex, columnIndex, err := grid.GetCoordinates(data.input)
			if data.shouldBeAnError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, data.rowIndex, rowIndex)
				assert.Equal(t, data.columnIndex, columnIndex)
			}
		})
	}
}
