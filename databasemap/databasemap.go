package databasemap

import "errors"

type DatabaseCell struct {
	ID       string
	Resident string
}

type DatabaseMap struct {
	ID    string
	Year  int
	Cells map[string]DatabaseCell
}

func NewBlankDatabaseMap() DatabaseMap {
	return DatabaseMap{
		Cells: make(map[string]DatabaseCell),
	}
}

func NewDatabaseMapWithTerritories(territoryNames []string) DatabaseMap {
	databaseMap := NewBlankDatabaseMap()
	for _, territoryName := range territoryNames {
		databaseMap.Cells[territoryName] = DatabaseCell{territoryName, ""}
	}
	return databaseMap
}

func (databaseMap *DatabaseMap) SetResident(territoryName string, nationID string) error {

	_, doesCellExist := databaseMap.Cells[territoryName]
	if !doesCellExist {
		return errors.New("Territory doesn't exist")
	}

	databaseMap.Cells[territoryName] = DatabaseCell{territoryName, nationID}

	return nil
}

func (databaseMap DatabaseMap) GetResident(territoryName string) (string, error) {

	territory, doesCellExist := databaseMap.Cells[territoryName]
	if !doesCellExist {
		return "", errors.New("Territory doesn't exist")
	}

	return territory.Resident, nil
}

func (databaseMap DatabaseMap) HasResident(territoryName string) (bool, error) {

	resident, err := databaseMap.GetResident(territoryName)
	if err != nil {
		return false, err
	}

	return len(resident) != 0, nil
}
