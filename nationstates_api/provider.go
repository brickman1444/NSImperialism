package nationstates_api

import "errors"

type NationStatesProvider interface {
	GetNationData(nationName string) (*Nation, error)
}

type NationStatesProviderSimpleMap struct {
	Nations map[string]Nation
}

func NewNationStatesProviderSimpleMap() NationStatesProviderSimpleMap {
	return NationStatesProviderSimpleMap{
		Nations: make(map[string]Nation),
	}
}

func (provider NationStatesProviderSimpleMap) GetNationData(nationName string) (*Nation, error) {
	nation, doesNationExist := provider.Nations[nationName]
	if !doesNationExist {
		return nil, errors.New("Nation doesn't exist")
	}

	return &nation, nil
}

func (provider *NationStatesProviderSimpleMap) PutNationData(nation Nation) {
	provider.Nations[nation.Id] = nation
}

var simpleMapInterfaceChecker NationStatesProvider = NationStatesProviderSimpleMap{}

type NationStatesProviderAPI struct {
}

func (provider NationStatesProviderAPI) GetNationData(nationName string) (*Nation, error) {
	return GetNationData(nationName)
}

var apiInterfaceChecker NationStatesProvider = NationStatesProviderAPI{}
