package main

import (
	"html/template"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/brickman1444/NSImperialism/grid"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/war"
	"github.com/stretchr/testify/assert"
)

func TestParseHTMLTemplates(t *testing.T) {

	allFiles, err := ioutil.ReadDir(".")
	assert.NoError(t, err)

	htmlFileNames := make([]string, 0)
	for _, fileInfo := range allFiles {
		if strings.HasSuffix(fileInfo.Name(), ".html") {
			htmlFileNames = append(htmlFileNames, fileInfo.Name())
		}
	}
	assert.NotEmpty(t, htmlFileNames)

	for _, htmlFileName := range htmlFileNames {
		t.Run(htmlFileName, func(t *testing.T) {
			_, err := template.ParseFiles(htmlFileName)
			assert.NoError(t, err)
		})
	}
}

func TestACompletedWarChangesOwnershipOfTheGrid(t *testing.T) {
	defender := &nationstates_api.Nation{}
	defender.SetDefenseForces(100)

	attacker := &nationstates_api.Nation{}
	attacker.SetDefenseForces(0)

	theWar := war.NewWar(attacker, defender, "", 1, 1)

	grid := grid.Grid{}
	grid.Rows[1].Cells[1].ResidentNation = defender

	wars := []*war.War{&theWar}

	tick(&grid, wars)

	assert.False(t, theWar.IsOngoing)
	assert.Equal(t, 100, theWar.Score)
	assert.Same(t, attacker, theWar.Advantage())

	assert.Same(t, attacker, grid.Rows[1].Cells[1].ResidentNation)
}
