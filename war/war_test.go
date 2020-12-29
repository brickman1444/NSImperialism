package war

import (
	"testing"

	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/stretchr/testify/assert"
)

func TestWarScoreChangePerYear(t *testing.T) {

	defender := &nationstates_api.Nation{}
	defender.SetDefenseForces(50)

	attacker := &nationstates_api.Nation{}
	attacker.SetDefenseForces(0)

	war := &War{attacker, defender, 0, ""}

	assert.Equal(t, 50, war.ScoreChangePerYear())
}

func TestNoOneHasAdvantageWhenWarScoreIsZero(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := &War{attacker, defender, 0, ""}

	assert.Nil(t, war.Advantage())
}

func TestAttackerHasAdvantageWhenWarHasPositiveScore(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := &War{attacker, defender, 1, ""}

	assert.Same(t, attacker, war.Advantage())
}

func TestDefenderHasAdvantageWhenWarHasNegativeScore(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := &War{attacker, defender, -1, ""}

	assert.Same(t, defender, war.Advantage())
}

func TestNoOneHasAdvantageWhenScoreIsZero(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	assert.Nil(t, Advantage(attacker, defender, 0))
}

func TestAttackerHasAdvantageWhenScoreIsPositive(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	assert.Same(t, attacker, Advantage(attacker, defender, 1))
}

func TestDefenderHasAdvantageWhenScoreIsNegative(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	assert.Same(t, defender, Advantage(attacker, defender, -1))
}
