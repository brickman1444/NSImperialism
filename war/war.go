package war

import (
	"fmt"
	"html/template"

	"github.com/brickman1444/NSImperialism/nationstates_api"
)

type War struct {
	Attacker *nationstates_api.Nation
	Defender *nationstates_api.Nation
	Score    int // 100 is attacker wins, -100 is defender wins
	Name     string
}

func (war *War) ScoreChangePerYear() int {
	return (100 - war.Attacker.GetDefenseForces()) - (100 - war.Defender.GetDefenseForces())
}

func (war *War) Advantage() *nationstates_api.Nation {
	if war.Score > 0 {
		return war.Attacker
	}

	if war.Score < 0 {
		return war.Defender
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
	return template.HTML(fmt.Sprintf("%d%%%s", absoluteScore, advantageDescription))
}
