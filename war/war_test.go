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

	war := NewWar(attacker, defender, "", "")

	assert.Equal(t, 50, war.ScoreChangePerYear())
}

func TestNoOneHasAdvantageWhenWarScoreIsZero(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := NewWar(attacker, defender, "", "")

	assert.Nil(t, war.Advantage())
}

func TestAttackerHasAdvantageWhenWarHasPositiveScore(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := NewWar(attacker, defender, "", "")
	war.Score = 1

	assert.Same(t, attacker, war.Advantage())
}

func TestDefenderHasAdvantageWhenWarHasNegativeScore(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := NewWar(attacker, defender, "", "")
	war.Score = -1

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

func TestNewWarIsOngoing(t *testing.T) {

	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := NewWar(attacker, defender, "", "")

	assert.True(t, war.IsOngoing)
}

func TestATickedWarChangesScore(t *testing.T) {

	defender := &nationstates_api.Nation{}
	defender.SetDefenseForces(50)

	attacker := &nationstates_api.Nation{}
	attacker.SetDefenseForces(0)

	war := NewWar(attacker, defender, "", "")

	assert.Equal(t, 0, war.Score)

	war.Tick()

	assert.Equal(t, 50, war.Score)

	war.Tick()

	assert.Equal(t, 100, war.Score)
}

func TestATickedWarCanEnd(t *testing.T) {

	defender := &nationstates_api.Nation{}
	defender.SetDefenseForces(50)

	attacker := &nationstates_api.Nation{}
	attacker.SetDefenseForces(0)

	war := NewWar(attacker, defender, "", "")

	assert.True(t, war.IsOngoing)

	didFinish := war.Tick()

	assert.False(t, didFinish)
	assert.True(t, war.IsOngoing)

	didFinish = war.Tick()

	assert.True(t, didFinish)
	assert.False(t, war.IsOngoing)
}

func TestFindOngoingWarFindsAWar(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	warAtA := NewWar(attacker, defender, "warAtA", "A")
	warAtB := NewWar(attacker, defender, "warAtB", "A")

	foundWar := FindOngoingWarAt([]War{warAtA, warAtB}, "A")

	assert.Equal(t, "warAtA", foundWar.Name)
}

func TestFindOngoingWarDoesntReturnACompletedWar(t *testing.T) {
	defender := &nationstates_api.Nation{}
	attacker := &nationstates_api.Nation{}

	war := NewWar(attacker, defender, "", "A")
	war.IsOngoing = false

	foundWar := FindOngoingWarAt([]War{war}, "A")

	assert.Nil(t, foundWar)
}
