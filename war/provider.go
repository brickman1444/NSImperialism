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

		if len(databaseWar.Attacker) == 0 {
			return nil, errors.New("Null attacker")
		}

		if len(databaseWar.Defender) == 0 {
			return nil, errors.New("Null defender")
		}

		retrievedWar := NewWar(databaseWar.Attacker, databaseWar.Defender, databaseWar.ID, databaseWar.TerritoryName)
		retrievedWar.IsOngoing = databaseWar.IsOngoing
		retrievedWar.Score = databaseWar.Score

		warsToReturn = append(warsToReturn, retrievedWar)
	}

	return warsToReturn, nil
}

func DatabaseWarFromRuntimeWar(war War) databasemap.DatabaseWar {
	return databasemap.DatabaseWar{
		Attacker:      war.AttackerID,
		Defender:      war.DefenderID,
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
