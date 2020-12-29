package grid

import (
	"testing"

	"github.com/brickman1444/NSImperialism/nationstates_api"
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

func TestColonizingAnEmptyCellPutsTheColonizerInTheCell(t *testing.T) {
	grid := Grid{}

	assert.Nil(t, grid.Rows[1].Cells[1].ResidentNation)

	nation := nationstates_api.Nation{}

	grid.Colonize(nation, "A1")

	assert.Equal(t, &nation, grid.Rows[1].Cells[1].ResidentNation)
}

func TestColonizingAnInvalidTargetProducesAnError(t *testing.T) {
	grid := Grid{}
	nation := nationstates_api.Nation{}

	err := grid.Colonize(nation, "A0")

	assert.Error(t, err)
}

func TestColonizingACellWithAResidentProducesAnErrorAndDoesntChangeTheGrid(t *testing.T) {
	grid := Grid{}
	firstNation := nationstates_api.Nation{}

	err := grid.Colonize(firstNation, "A1")
	assert.NoError(t, err)

	secondNation := nationstates_api.Nation{}
	err = grid.Colonize(secondNation, "A1")
	assert.Error(t, err)
	assert.Equal(t, &firstNation, grid.Rows[1].Cells[1].ResidentNation)
}
