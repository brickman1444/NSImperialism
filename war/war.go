package war

import (
	"fmt"
	"html/template"
	"math/rand"

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

func NewWar(attacker *nationstates_api.Nation, defender *nationstates_api.Nation, name string, territoryName string) War {
	return War{attacker.Id, defender.Id, 0, name, territoryName, true}
}

func (war *War) Advantage() (*string, error) {
	return Advantage(war.AttackerID, war.DefenderID, war.Score), nil
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

	advantageID, err := war.Advantage()
	if err != nil {
		return "", err
	}

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

func (war *War) Tick(nationStatesProvider nationstates_api.NationStatesProvider) (bool, error) {

	if war.IsOngoing {

		defender, err := nationStatesProvider.GetNationData(war.DefenderID)
		if err != nil {
			return false, err
		}

		attacker, err := nationStatesProvider.GetNationData(war.AttackerID)
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
