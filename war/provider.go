package war

import (
	"errors"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/nationstates_api"
)

func GetWars(databaseMap databasemap.DatabaseMap, nationStatesProvider nationstates_api.NationStatesProvider) ([]War, error) {

	databaseWars := databaseMap.GetWars()

	warsToReturn := make([]War, 0, len(databaseWars))
	for _, databaseWar := range databaseWars {

		attacker, err := nationStatesProvider.GetNationData(databaseWar.Attacker)
		if err != nil {
			return nil, err
		}

		if attacker == nil {
			return nil, errors.New("Null attacker")
		}

		defender, err := nationStatesProvider.GetNationData(databaseWar.Defender)
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

func DatabaseWarFromRuntimeWar(war War) databasemap.DatabaseWar {
	return databasemap.DatabaseWar{
		Attacker:      war.Attacker.Id,
		Defender:      war.Defender.Id,
		Score:         war.Score,
		ID:            war.Name,
		TerritoryName: war.TerritoryName,
		IsOngoing:     war.IsOngoing,
	}
}

func PutWars(databaseMap *databasemap.DatabaseMap, warsToAdd []War) {

	databaseWarsToPut := make([]databasemap.DatabaseWar, 0, len(warsToAdd))
	for _, warToAdd := range warsToAdd {
		databaseWarsToPut = append(databaseWarsToPut, DatabaseWarFromRuntimeWar(warToAdd))
	}

	databaseMap.PutWars(databaseWarsToPut)
}
