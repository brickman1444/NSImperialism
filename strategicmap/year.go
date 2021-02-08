package strategicmap

import "github.com/brickman1444/NSImperialism/dynamodbwrapper"

type YearInterface interface {
	Get() (int, error)
	Increment() error
}

type YearSimpleProvider struct {
	Year int
}

func (simpleProvider YearSimpleProvider) Get() (int, error) {
	return simpleProvider.Year, nil
}

func (simpleProvider *YearSimpleProvider) Increment() error {
	simpleProvider.Year++
	return nil
}

var simpleYearInterfaceChecker YearInterface = &YearSimpleProvider{}

type YearDatabaseProvider struct {
}

func (databaseProvider YearDatabaseProvider) Get() (int, error) {

	databaseMap, err := dynamodbwrapper.GetMap()
	if err != nil {
		return 0, err
	}

	return databaseMap.Year, nil
}

func (databaseProvider *YearDatabaseProvider) Increment() error {

	databaseMap, err := dynamodbwrapper.GetMap()
	if err != nil {
		return err
	}

	databaseMap.Year++

	return dynamodbwrapper.PutMap(databaseMap)
}

var databaseYearInterfaceChecker YearInterface = &YearDatabaseProvider{}
