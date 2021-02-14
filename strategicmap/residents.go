package strategicmap

import "github.com/brickman1444/NSImperialism/dynamodbwrapper"

type ResidentsInterface interface {
	SetResident(territoryName string, nationID string) error
	GetResident(territoryName string) (string, error)
	HasResident(territoryName string) (bool, error)
	CanExpand(nationID string) (bool, error)
	GetAllMapIDs() ([]string, error)
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

func (simpleMap ResidentsSimpleMap) CanExpand(nationID string) (bool, error) {
	return true, nil
}

func (simpleMap ResidentsSimpleMap) GetAllMapIDs() ([]string, error) {
	return []string{}, nil
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

func (database ResidentsDatabase) CanExpand(nationID string) (bool, error) {
	return true, nil
}

func (database ResidentsDatabase) GetAllMapIDs() ([]string, error) {
	return dynamodbwrapper.GetAllMapIDs()
}

var databaseInterfaceChecker ResidentsInterface = ResidentsDatabase{}
