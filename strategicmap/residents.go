package strategicmap

import (
	"errors"
	"math/rand"

	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
)

type ResidentsInterface interface {
	SetResident(territoryName string, nationID string) error
	GetResident(territoryName string) (string, error)
	HasResident(territoryName string) (bool, error)
}

type ResidentsSimpleMap struct {
	residents map[string]string
}

func NewResidentsSimpleMap() ResidentsSimpleMap {
	return ResidentsSimpleMap{residents: make(map[string]string)}
}

func (simpleMap ResidentsSimpleMap) SetResident(territoryName string, nationID string) error {
	simpleMap.residents[territoryName] = nationID
	return nil
}

func (simpleMap ResidentsSimpleMap) GetResident(territoryName string) (string, error) {
	return simpleMap.residents[territoryName], nil
}

func (simpleMap ResidentsSimpleMap) HasResident(territoryName string) (bool, error) {
	_, doesExist := simpleMap.residents[territoryName]
	return doesExist, nil
}

var simpleMapInterfaceChecker ResidentsInterface = ResidentsSimpleMap{}

type ResidentsDatabase struct {
}

func (database ResidentsDatabase) SetResident(territoryName string, nationID string) error {
	return dynamodbwrapper.PutCell(dynamodbwrapper.DatabaseCell{ID: territoryName, Resident: nationID})
}

func (database ResidentsDatabase) GetResident(territoryName string) (string, error) {

	databaseCell, err := dynamodbwrapper.GetCell(territoryName)
	if err == dynamodbwrapper.CellDoesntExistError {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return databaseCell.Resident, nil
}

func (database ResidentsDatabase) HasResident(territoryName string) (bool, error) {

	_, err := dynamodbwrapper.GetCell(territoryName)
	if err == dynamodbwrapper.CellDoesntExistError {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func MakeNewRandomMap(mapLayout Map, participatingNations []string) (dynamodbwrapper.DatabaseMap, error) {
	databaseMap := dynamodbwrapper.NewDatabaseMap()

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
		databaseMap.Cells[territoryID] = dynamodbwrapper.DatabaseCell{
			ID:       territoryID,
			Resident: residentsForEachCell[territoryIndex],
		}
	}

	return databaseMap, nil
}

var databaseInterfaceChecker ResidentsInterface = ResidentsDatabase{}
