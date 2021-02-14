package strategicmap

import "github.com/brickman1444/NSImperialism/dynamodbwrapper"

type YearInterface interface {
	Get(mapID string) (int, error)
	Increment(mapID string) error
}

type YearSimpleProvider struct {
	Year int
}

func (simpleProvider YearSimpleProvider) Get(mapID string) (int, error) {
	return simpleProvider.Year, nil
}

func (simpleProvider *YearSimpleProvider) Increment(mapID string) error {
	simpleProvider.Year++
	return nil
}

var simpleYearInterfaceChecker YearInterface = &YearSimpleProvider{}

type YearDatabaseProvider struct {
}

func (databaseProvider YearDatabaseProvider) Get(mapID string) (int, error) {

	databaseMap, err := dynamodbwrapper.GetMap(mapID)
	if err != nil {
		return 0, err
	}

	return databaseMap.Year, nil
}

func (databaseProvider *YearDatabaseProvider) Increment(mapID string) error {

	databaseMap, err := dynamodbwrapper.GetMap(mapID)
	if err != nil {
		return err
	}

	databaseMap.Year++

	return dynamodbwrapper.PutMap(databaseMap)
}

var databaseYearInterfaceChecker YearInterface = &YearDatabaseProvider{}
