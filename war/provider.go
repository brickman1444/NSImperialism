package war

import (
	"github.com/brickman1444/NSImperialism/databasemap"
)

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

func PutWars(databaseMap *databasemap.DatabaseMap, warsToAdd []databasemap.DatabaseWar) {

	databaseMap.PutWars(warsToAdd)
}
