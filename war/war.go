package war

import (
	"fmt"
	"html/template"

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

func NewWar(attacker *nationstates_api.Nation, defender *nationstates_api.Nation, name string, territoryName string) War {
	return War{attacker, defender, 0, name, territoryName, true}
}

func (war *War) ScoreChangePerYear() int {
	return (100 - war.Attacker.GetDefenseForces()) - (100 - war.Defender.GetDefenseForces())
}

func (war *War) Advantage() *nationstates_api.Nation {
	return Advantage(war.Attacker, war.Defender, war.Score)
}

func Advantage(attacker *nationstates_api.Nation, defender *nationstates_api.Nation, score int) *nationstates_api.Nation {
	if score > 0 {
		return attacker
	}

	if score < 0 {
		return defender
	}

	return nil
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (war *War) ScoreDescription() template.HTML {

	advantage := war.Advantage()
	advantageDescription := ""
	if advantage != nil {
		advantageDescription = fmt.Sprintf(" in favor of %s", string(advantage.FlagAndName()))
	}

	absoluteScore := Abs(war.Score)
	return template.HTML(fmt.Sprintf("Currently %d%%%s", absoluteScore, advantageDescription))
}

func (war *War) ScorePerYearDescription() template.HTML {

	advantage := Advantage(war.Attacker, war.Defender, war.ScoreChangePerYear())
	advantageDescription := ""
	if advantage != nil {
		advantageDescription = fmt.Sprintf(" in favor of %s", string(advantage.FlagAndName()))
	}

	absoluteScore := Abs(war.ScoreChangePerYear())
	return template.HTML(fmt.Sprintf("+%d%%%s per year", absoluteScore, advantageDescription))
}

func FindOngoingWarAt(wars []*War, territoryName string) *War {
	for warIndex, war := range wars {
		if war.TerritoryName == territoryName && war.IsOngoing {
			return wars[warIndex]
		}
	}
	return nil
}

func (war *War) Tick() bool {

	if war.IsOngoing {
		war.Score += war.ScoreChangePerYear()

		if Abs(war.Score) >= 100 {
			war.IsOngoing = false
			return true
		}
	}

	return false
}
