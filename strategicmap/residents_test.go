package strategicmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomMapHasEveryCellFilled(t *testing.T) {

	staticMap := Map{Territories: []Territory{
		{"A", 0, 0},
		{"B", 0, 0},
		{"C", 0, 0},
	}}

	for simulationIndex := 0; simulationIndex < 1000; simulationIndex++ {
		randomMap, err := MakeNewRandomMap(staticMap, []string{"nation1", "nation2"})
		assert.NoError(t, err)

		assert.NotEmpty(t, randomMap.Cells["A"].Resident)
		assert.NotEmpty(t, randomMap.Cells["B"].Resident)
		assert.NotEmpty(t, randomMap.Cells["C"].Resident)
	}
}

func TestRandomMapEachNationHasAtLeastOneCell(t *testing.T) {

	staticMap := Map{Territories: []Territory{
		{"A", 0, 0},
		{"B", 0, 0},
		{"C", 0, 0},
	}}

	for simulationIndex := 0; simulationIndex < 1000; simulationIndex++ {
		randomMap, err := MakeNewRandomMap(staticMap, []string{"nation1", "nation2"})
		assert.NoError(t, err)

		numberOfNation1Cells := 0
		numberOfNation2Cells := 0

		for _, cell := range randomMap.Cells {
			if cell.Resident == "nation1" {
				numberOfNation1Cells++
			}

			if cell.Resident == "nation2" {
				numberOfNation2Cells++
			}
		}

		assert.GreaterOrEqual(t, numberOfNation1Cells, 1)
		assert.LessOrEqual(t, numberOfNation1Cells, 2)

		assert.GreaterOrEqual(t, numberOfNation2Cells, 1)
		assert.LessOrEqual(t, numberOfNation2Cells, 2)
	}
}

func TestCreatingAMapWithLessThanTwoNationsIsAnError(t *testing.T) {

	staticMap := Map{Territories: []Territory{
		{"A", 0, 0},
		{"B", 0, 0},
		{"C", 0, 0},
	}}

	for simulationIndex := 0; simulationIndex < 1000; simulationIndex++ {
		_, err := MakeNewRandomMap(staticMap, []string{"nation1"})
		assert.Error(t, err)
	}
}

func TestCreatingAMapMoreNationsThanCellsIsAnError(t *testing.T) {

	staticMap := Map{Territories: []Territory{
		{"A", 0, 0},
		{"B", 0, 0},
		{"C", 0, 0},
	}}

	for simulationIndex := 0; simulationIndex < 1000; simulationIndex++ {
		_, err := MakeNewRandomMap(staticMap, []string{"nation1", "nation2", "nation3", "nation4"})
		assert.Error(t, err)
	}
}

func TestRandomMapHasNonEmptyID(t *testing.T) {

	staticMap := Map{Territories: []Territory{
		{"A", 0, 0},
		{"B", 0, 0},
		{"C", 0, 0},
	}}

	databaseMap, err := MakeNewRandomMap(staticMap, []string{"nation1", "nation2"})
	assert.NoError(t, err)
	assert.NotEmpty(t, databaseMap.ID)
}
