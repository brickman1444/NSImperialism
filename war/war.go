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

func ScoreDescription(war databasemap.DatabaseWar, attacker nationstates_api.Nation, defender nationstates_api.Nation) template.HTML {

	advantageID := WarAdvantage(war)

	advantageDescription := ""
	if advantageID != nil {
		if *advantageID == attacker.Id {
			advantageDescription = fmt.Sprintf(" in favor of %s", string(attacker.FlagAndName()))
		} else if *advantageID == defender.Id {
			advantageDescription = fmt.Sprintf(" in favor of %s", string(defender.FlagAndName()))
		}
	}

	absoluteScore := Abs(war.Score)
	return template.HTML(fmt.Sprintf("Currently %d%%%s", absoluteScore, advantageDescription))
}

func FindOngoingWarAt(wars []databasemap.DatabaseWar, territoryName string) *databasemap.DatabaseWar {
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

type RenderedWar struct {
	IsOngoing        bool
	Name             string
	Attacker         template.HTML
	Defender         template.HTML
	ScoreDescription template.HTML
}

func RenderWar(war databasemap.DatabaseWar, nationStatesProvider nationstates_api.NationStatesProvider) (RenderedWar, error) {

	attacker, err := nationStatesProvider.GetNationData(war.Attacker)
	if err != nil {
		return RenderedWar{}, err
	}

	defender, err := nationStatesProvider.GetNationData(war.Defender)
	if err != nil {
		return RenderedWar{}, err
	}

	return RenderedWar{
		IsOngoing:        war.IsOngoing,
		Name:             war.ID,
		Attacker:         attacker.FlagAndName(),
		Defender:         defender.FlagAndName(),
		ScoreDescription: ScoreDescription(war, *attacker, *defender),
	}, nil
}

func RenderWars(wars []databasemap.DatabaseWar, nationStatesProvider nationstates_api.NationStatesProvider) ([]RenderedWar, error) {
	renderedWars := []RenderedWar{}
	for _, war := range wars {
		renderedWar, err := RenderWar(war, nationStatesProvider)
		if err != nil {
			return []RenderedWar{}, err
		}

		renderedWars = append(renderedWars, renderedWar)
	}
	return renderedWars, nil
}
