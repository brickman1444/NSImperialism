package war

import (
	"math"
	"testing"

	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/stretchr/testify/assert"
)

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

	scoreTurnZero := war.Score
	assert.Equal(t, 0, scoreTurnZero)

	war.Tick()

	scoreTurnOne := war.Score
	assert.NotEqual(t, scoreTurnZero, scoreTurnOne)

	war.Tick()

	scoreTurnTwo := war.Score
	assert.NotEqual(t, scoreTurnOne, scoreTurnTwo)
}

func TestATickedWarCanEnd(t *testing.T) {

	defender := &nationstates_api.Nation{}
	defender.SetDefenseForces(50)

	attacker := &nationstates_api.Nation{}
	attacker.SetDefenseForces(0)

	war := NewWar(attacker, defender, "", "")

	assert.True(t, war.IsOngoing)

	turnCount := 0
	maximumTurnCount := 1000
	finalTickResult := false
	for war.IsOngoing && turnCount < maximumTurnCount {
		finalTickResult = war.Tick()
		turnCount++
	}

	assert.True(t, finalTickResult)
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

func TestMorePowerfulNationDoesntAlwaysWinWar(t *testing.T) {

	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(60)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(40)

	attackerWinCount := 0
	defenderWinCount := 0
	totalLength := 0
	minimumLength := math.MaxInt32
	maximumLength := 0

	totalNumberOfSimulations := 10000

	for warIndex := 0; warIndex < totalNumberOfSimulations; warIndex++ {
		war := NewWar(attacker, defender, "", "")

		length := 0
		for war.IsOngoing {
			war.Tick()
			length++
		}

		winner := war.Advantage()
		if winner != nil && winner.Id == attacker.Id {
			attackerWinCount++
		} else {
			defenderWinCount++
		}

		totalLength = totalLength + length

		if length < minimumLength {
			minimumLength = length
		}
		if length > maximumLength {
			maximumLength = length
		}
	}

	averageLength := float32(totalLength) / float32(totalNumberOfSimulations)

	assert.Greater(t, attackerWinCount, totalNumberOfSimulations/10)
	assert.Greater(t, defenderWinCount, totalNumberOfSimulations/10)
	assert.Greater(t, attackerWinCount, defenderWinCount)
	assert.Greater(t, minimumLength, 2)
	assert.Less(t, maximumLength, 70)
	assert.Less(t, averageLength, float32(9))
}
