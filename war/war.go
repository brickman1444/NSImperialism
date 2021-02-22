package war

import (
	"fmt"
	"html/template"
	"math/rand"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/nationstates_api"
)

type War struct {
	AttackerID    string
	DefenderID    string
	Score         int // 100 is attacker wins, -100 is defender wins
	Name          string
	TerritoryName string
	IsOngoing     bool
}

func NewWar(attackerID string, defenderID string, name string, territoryName string) War {
	return War{attackerID, defenderID, 0, name, territoryName, true}
}

func WarAdvantage(war databasemap.DatabaseWar) *string {
	return Advantage(war.Attacker, war.Defender, war.Score)
}

func Advantage(attackerID string, defenderID string, score int) *string {
	if score > 0 {
		return &attackerID
	}

	if score < 0 {
		return &defenderID
	}

	return nil
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (war *War) ScoreDescription(nationStatesProvider nationstates_api.NationStatesProvider) (template.HTML, error) {

	advantageID := WarAdvantage(DatabaseWarFromRuntimeWar(*war))

	advantageDescription := ""
	if advantageID != nil {

		advantageNation, err := nationStatesProvider.GetNationData(*advantageID)
		if err != nil {
			return "", err
		}

		advantageDescription = fmt.Sprintf(" in favor of %s", string(advantageNation.FlagAndName()))
	}

	absoluteScore := Abs(war.Score)
	return template.HTML(fmt.Sprintf("Currently %d%%%s", absoluteScore, advantageDescription)), nil
}

func FindOngoingWarAt(wars []War, territoryName string) *War {
	for warIndex, war := range wars {
		if war.TerritoryName == territoryName && war.IsOngoing {
			return &wars[warIndex]
		}
	}
	return nil
}

const battleScoreDelta = 40 // TODO: This should be 10 * number of years the war has been ongoing

func Tick(war *databasemap.DatabaseWar, nationStatesProvider nationstates_api.NationStatesProvider) (bool, error) {

	if war.IsOngoing {

		defender, err := nationStatesProvider.GetNationData(war.Defender)
		if err != nil {
			return false, err
		}

		attacker, err := nationStatesProvider.GetNationData(war.Attacker)
		if err != nil {
			return false, err
		}

		// TODO: This should be divided by the number of territories controlled
		defenderDefenseForcesInverted := 100 - defender.GetDefenseForces()
		attackerDefenseForcesInverted := 100 - attacker.GetDefenseForces()

		randomRoll := rand.Intn(defenderDefenseForcesInverted + attackerDefenseForcesInverted)

		if randomRoll < defenderDefenseForcesInverted {
			war.Score -= battleScoreDelta
		} else {
			war.Score += battleScoreDelta
		}

		if Abs(war.Score) >= 100 {
			war.IsOngoing = false
			return true, nil
		}
	}

	return false, nil
}
