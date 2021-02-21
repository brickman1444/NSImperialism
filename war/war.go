package war

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"

	"github.com/brickman1444/NSImperialism/nationstates_api"
)

type War struct {
	Attacker      *nationstates_api.Nation
	Defender      *nationstates_api.Nation
	Score         int // 100 is attacker wins, -100 is defender wins
	Name          string
	TerritoryName string
	IsOngoing     bool
}

func (war War) AttackerNation(nationStatesProvider nationstates_api.NationStatesProvider) (nationstates_api.Nation, error) {
	if war.Attacker == nil {
		return nationstates_api.Nation{}, errors.New("nil nation")
	} else {
		return *war.Attacker, nil
	}
}

func (war War) DefenderNation(nationStatesProvider nationstates_api.NationStatesProvider) (nationstates_api.Nation, error) {
	if war.Defender == nil {
		return nationstates_api.Nation{}, errors.New("nil nation")
	} else {
		return *war.Defender, nil
	}
}

func NewWar(attacker *nationstates_api.Nation, defender *nationstates_api.Nation, name string, territoryName string) War {
	return War{attacker, defender, 0, name, territoryName, true}
}

func (war *War) Advantage(nationStatesProvider nationstates_api.NationStatesProvider) (*string, error) {

	attacker, err := war.AttackerNation(nationStatesProvider)
	if err != nil {
		return nil, err
	}

	defender, err := war.DefenderNation(nationStatesProvider)
	if err != nil {
		return nil, err
	}

	return Advantage(attacker.Id, defender.Id, war.Score), nil
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

	advantageID, err := war.Advantage(nationStatesProvider)
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

func (war *War) Tick() bool {

	if war.IsOngoing {

		// TODO: This should be divided by the number of territories controlled
		defenderDefenseForcesInverted := 100 - war.Defender.GetDefenseForces()
		attackerDefenseForcesInverted := 100 - war.Attacker.GetDefenseForces()

		randomRoll := rand.Intn(defenderDefenseForcesInverted + attackerDefenseForcesInverted)

		if randomRoll < defenderDefenseForcesInverted {
			war.Score -= battleScoreDelta
		} else {
			war.Score += battleScoreDelta
		}

		if Abs(war.Score) >= 100 {
			war.IsOngoing = false
			return true
		}
	}

	return false
}
