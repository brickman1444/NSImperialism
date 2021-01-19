package war

import (
	"errors"

	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
	"github.com/brickman1444/NSImperialism/nationstates_api"
)

type WarProviderInterface interface {
	GetWars() ([]War, error)
	PutWars(wars []War) error
}

type WarProviderSimpleList struct {
	wars []War
}

func NewWarProviderSimpleList() WarProviderSimpleList {
	return WarProviderSimpleList{}
}

func (simpleList WarProviderSimpleList) GetWars() ([]War, error) {
	return simpleList.wars, nil
}

func (simpleList *WarProviderSimpleList) PutWars(warsToAdd []War) error {

	for _, warToAdd := range warsToAdd {
		foundExistingWar := false
		for existingWarIndex, existingWar := range simpleList.wars {
			if existingWar.Name == warToAdd.Name {
				simpleList.wars[existingWarIndex] = warToAdd
				foundExistingWar = true
				break
			}
		}
		if !foundExistingWar {
			simpleList.wars = append(simpleList.wars, warToAdd)
		}
	}

	return nil
}

var warProviderSimpleListInterfaceChecker WarProviderInterface = &WarProviderSimpleList{}

type WarProviderDatabase struct {
}

func (provider WarProviderDatabase) GetWars() ([]War, error) {

	databaseWars, err := dynamodbwrapper.GetWars()
	if err != nil {
		return nil, err
	}

	warsToReturn := make([]War, 0, len(databaseWars))
	for _, databaseWar := range databaseWars {

		attacker, err := nationstates_api.GetNationData(databaseWar.Attacker)
		if err != nil {
			return nil, err
		}

		if attacker == nil {
			return nil, errors.New("Null attacker")
		}

		defender, err := nationstates_api.GetNationData(databaseWar.Defender)
		if err != nil {
			return nil, err
		}

		if defender == nil {
			return nil, errors.New("Null defender")
		}

		retrievedWar := NewWar(attacker, defender, databaseWar.ID, databaseWar.TerritoryName)
		retrievedWar.IsOngoing = databaseWar.IsOngoing
		retrievedWar.Score = databaseWar.Score

		warsToReturn = append(warsToReturn, retrievedWar)
	}

	return warsToReturn, nil
}

func DatabaseWarFromRuntimeWar(war War) dynamodbwrapper.DatabaseWar {
	return dynamodbwrapper.DatabaseWar{
		Attacker:      war.Attacker.Id,
		Defender:      war.Defender.Id,
		Score:         war.Score,
		ID:            war.Name,
		TerritoryName: war.TerritoryName,
		IsOngoing:     war.IsOngoing,
	}
}

func (provider WarProviderDatabase) PutWars(warsToAdd []War) error {

	databaseWarsToPut := make([]dynamodbwrapper.DatabaseWar, 0, len(warsToAdd))
	for _, warToAdd := range warsToAdd {
		databaseWarsToPut = append(databaseWarsToPut, DatabaseWarFromRuntimeWar(warToAdd))
	}

	return dynamodbwrapper.PutWars(databaseWarsToPut)
}

var warProviderDatabaseInterfaceChecker WarProviderInterface = WarProviderDatabase{}
