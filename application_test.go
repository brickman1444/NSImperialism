package main

import (
	"html/template"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/brickman1444/NSImperialism/databasemap"
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

func TestACompletedWarChangesResidenceOfTheTerritory(t *testing.T) {

	for simulationCount := 0; simulationCount < 1000; simulationCount++ {
		defender := &nationstates_api.Nation{Id: "Defender"}
		defender.SetDefenseForces(100)

		attacker := &nationstates_api.Nation{Id: "Attacker"}
		attacker.SetDefenseForces(0)

		theWar := war.NewWar(attacker, defender, "", "A")

		residentNations := databasemap.NewDatabaseMapWithTerritories([]string{"A"})
		err := residentNations.SetResident("A", defender.Id)
		assert.NoError(t, err)

		warProvider := war.NewWarProviderSimpleList()
		err = warProvider.PutWars([]war.War{theWar})
		assert.NoError(t, err)

		for warTurnCount := 0; warTurnCount < 1000; warTurnCount++ {

			tick(&residentNations, &warProvider)

			wars, err := warProvider.GetWars()
			assert.NoError(t, err)
			assert.Len(t, wars, 1)

			if !wars[0].IsOngoing {
				break
			}
		}

		wars, err := warProvider.GetWars()
		assert.NoError(t, err)
		assert.Len(t, wars, 1)

		finishedWar := wars[0]
		assert.False(t, finishedWar.IsOngoing)

		if finishedWar.Advantage().Id == attacker.Id {
			newResidentID, err := residentNations.GetResident("A")

			assert.NoError(t, err)

			assert.Equal(t, attacker.Id, newResidentID)
		} else {
			newResidentID, err := residentNations.GetResident("A")

			assert.NoError(t, err)

			assert.Equal(t, defender.Id, newResidentID)
		}
	}
}

func TestApplicationTickUpdatesWars(t *testing.T) {
	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(100)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(0)

	residentNations := databasemap.NewDatabaseMapWithTerritories([]string{"A"})
	residentNations.SetResident("A", defender.Id)

	warProvider := war.NewWarProviderSimpleList()
	err := warProvider.PutWars([]war.War{war.NewWar(attacker, defender, "warForA", "A")})
	assert.NoError(t, err)

	tick(&residentNations, &warProvider)

	retrievedWars, err := warProvider.GetWars()
	assert.NoError(t, err)

	assert.Len(t, retrievedWars, 1)

	assert.Equal(t, "warForA", retrievedWars[0].Name)
	assert.NotEqual(t, 0, retrievedWars[0].Score)
}
