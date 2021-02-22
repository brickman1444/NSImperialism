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

		nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()

		defender := &nationstates_api.Nation{Id: "Defender"}
		defender.SetDefenseForces(100)
		nationStatesProvider.PutNationData(*defender)

		attacker := &nationstates_api.Nation{Id: "Attacker"}
		attacker.SetDefenseForces(0)
		nationStatesProvider.PutNationData(*attacker)

		theWar := databasemap.NewWar(attacker.Id, defender.Id, "", "A")

		residentNations := databasemap.NewDatabaseMapWithTerritories([]string{"A"})
		err := residentNations.SetResident("A", defender.Id)
		assert.NoError(t, err)

		war.PutWars(&residentNations, []databasemap.DatabaseWar{theWar})

		for warTurnCount := 0; warTurnCount < 1000; warTurnCount++ {

			tick(&residentNations, nationStatesProvider)

			wars := residentNations.GetWars()
			assert.Len(t, wars, 1)

			if !wars[0].IsOngoing {
				break
			}
		}

		wars := residentNations.GetWars()
		assert.Len(t, wars, 1)

		finishedWar := wars[0]
		assert.False(t, finishedWar.IsOngoing)

		advantageID := war.WarAdvantage(finishedWar)
		assert.NotNil(t, advantageID)

		if *advantageID == attacker.Id {
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

	nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()
	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(100)
	nationStatesProvider.PutNationData(*defender)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(0)
	nationStatesProvider.PutNationData(*attacker)

	residentNations := databasemap.NewDatabaseMapWithTerritories([]string{"A"})
	residentNations.SetResident("A", defender.Id)

	war.PutWars(&residentNations, []databasemap.DatabaseWar{databasemap.NewWar(attacker.Id, defender.Id, "warForA", "A")})

	tick(&residentNations, nationStatesProvider)

	retrievedWars := residentNations.GetWars()

	assert.Len(t, retrievedWars, 1)

	assert.Equal(t, "warForA", retrievedWars[0].ID)
	assert.NotEqual(t, 0, retrievedWars[0].Score)
}
