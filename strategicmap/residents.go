package strategicmap

type ResidentsInterface interface {
	SetResident(territoryName string, nationID string)
	GetResident(territoryName string) string
	HasResident(territoryName string) bool
}

type ResidentsSimpleMap struct {
	residents map[string]string
}

func NewResidentsSimpleMap() ResidentsSimpleMap {
	return ResidentsSimpleMap{residents: make(map[string]string)}
}

func (simpleMap ResidentsSimpleMap) SetResident(territoryName string, nationID string) {
	simpleMap.residents[territoryName] = nationID
}

func (simpleMap ResidentsSimpleMap) GetResident(territoryName string) string {
	return simpleMap.residents[territoryName]
}

func (simpleMap ResidentsSimpleMap) HasResident(territoryName string) bool {
	_, doesExist := simpleMap.residents[territoryName]
	return doesExist
}

var interfaceChecker ResidentsSimpleMap = ResidentsSimpleMap{}
