package main

import (
	"html/template"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/strategicmap"
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

func TestACompletedWarChangesResidenceOfTheTerritory(t *testing.T) {
	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(100)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(0)

	theWar := war.NewWar(attacker, defender, "", "A")

	residentNations := strategicmap.NewResidentsSimpleMap()
	residentNations.SetResident("A", defender.Id)

	wars := []*war.War{&theWar}
	year := 0

	tick(residentNations, wars, &year)

	assert.False(t, theWar.IsOngoing)
	assert.Equal(t, 100, theWar.Score)
	assert.Equal(t, attacker.Id, theWar.Advantage().Id)

	newResidentID, err := residentNations.GetResident("A")

	assert.NoError(t, err)

	assert.Equal(t, attacker.Id, newResidentID)
}
