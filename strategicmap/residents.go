package strategicmap

import (
	"errors"
	"math/rand"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
	"github.com/google/uuid"
)

type MapsInterface interface {
	GetMap(mapID string) (databasemap.DatabaseMap, error)
	PutMap(databaseMap databasemap.DatabaseMap) error
}

type MapsDatabase struct {
}

func (mapsDatabase MapsDatabase) GetMap(mapID string) (databasemap.DatabaseMap, error) {
	return dynamodbwrapper.GetMap(mapID)
}

func (mapsDatabase MapsDatabase) PutMap(databaseMap databasemap.DatabaseMap) error {
	return dynamodbwrapper.PutMap(databaseMap)
}

var databaseInterfaceChecker MapsInterface = MapsDatabase{}

func MakeNewRandomMap(mapLayout Map, participatingNations []string) (databasemap.DatabaseMap, error) {
	databaseMap := databasemap.NewBlankDatabaseMap()

	databaseMap.ID = uuid.NewString()

	if len(participatingNations) < 2 {
		return databaseMap, errors.New("Creating a map requires at least two nations")
	}

	if len(participatingNations) > len(mapLayout.Territories) {
		return databaseMap, errors.New("There must be space for each nation to get at least one territory")
	}

	residentsForEachCell := make([]string, len(participatingNations), len(mapLayout.Territories))
	copy(residentsForEachCell, participatingNations) // this gives each nation one cell

	for len(residentsForEachCell) < len(mapLayout.Territories) {

		randomNationIndex := rand.Intn(len(participatingNations))
		residentsForEachCell = append(residentsForEachCell, participatingNations[randomNationIndex])
	}

	rand.Shuffle(len(residentsForEachCell), func(i, j int) {
		residentsForEachCell[i], residentsForEachCell[j] = residentsForEachCell[j], residentsForEachCell[i]
	})

	for territoryIndex, _ := range mapLayout.Territories {
		territoryID := mapLayout.Territories[territoryIndex].Name
		databaseMap.Cells[territoryID] = databasemap.DatabaseCell{
			ID:       territoryID,
			Resident: residentsForEachCell[territoryIndex],
		}
	}

	return databaseMap, nil
}
