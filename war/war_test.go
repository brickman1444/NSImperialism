package war

import (
	"math"
	"testing"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/stretchr/testify/assert"
)

func TestNoOneHasAdvantageWhenWarScoreIsZero(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "")

	advantage := WarAdvantage(war)

	assert.Nil(t, advantage)
}

func TestAttackerHasAdvantageWhenWarHasPositiveScore(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "")
	war.Score = 1

	advantageID := WarAdvantage(war)
	assert.NotNil(t, advantageID)

	assert.Equal(t, "attacker", *advantageID)
}

func TestDefenderHasAdvantageWhenWarHasNegativeScore(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "")
	war.Score = -1

	advantageID := WarAdvantage(war)
	assert.NotNil(t, advantageID)

	assert.Equal(t, "defender", *advantageID)
}

func TestNoOneHasAdvantageWhenScoreIsZero(t *testing.T) {
	assert.Nil(t, Advantage("attacker", "defender", 0))
}

func TestAttackerHasAdvantageWhenScoreIsPositive(t *testing.T) {
	assert.Equal(t, "attacker", *Advantage("attacker", "defender", 1))
}

func TestDefenderHasAdvantageWhenScoreIsNegative(t *testing.T) {
	assert.Equal(t, "defender", *Advantage("attacker", "defender", -1))
}

func TestNewWarIsOngoing(t *testing.T) {
	war := NewWar("", "", "", "")

	assert.True(t, war.IsOngoing)
}

func TestATickedWarChangesScore(t *testing.T) {

	defender := nationstates_api.Nation{Id: "defender"}
	defender.SetDefenseForces(50)

	attacker := nationstates_api.Nation{Id: "attacker"}
	attacker.SetDefenseForces(0)

	nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()
	nationStatesProvider.PutNationData(defender)
	nationStatesProvider.PutNationData(attacker)

	war := databasemap.NewWar(attacker.Id, defender.Id, "", "")

	scoreTurnZero := war.Score
	assert.Equal(t, 0, scoreTurnZero)

	Tick(&war, nationStatesProvider)

	scoreTurnOne := war.Score
	assert.NotEqual(t, scoreTurnZero, scoreTurnOne)

	Tick(&war, nationStatesProvider)

	scoreTurnTwo := war.Score
	assert.NotEqual(t, scoreTurnOne, scoreTurnTwo)
}

func TestATickedWarCanEnd(t *testing.T) {

	defender := nationstates_api.Nation{Id: "defender"}
	defender.SetDefenseForces(50)

	attacker := nationstates_api.Nation{Id: "attacker"}
	attacker.SetDefenseForces(0)

	nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()
	nationStatesProvider.PutNationData(defender)
	nationStatesProvider.PutNationData(attacker)

	war := databasemap.NewWar(attacker.Id, defender.Id, "", "")

	assert.True(t, war.IsOngoing)

	turnCount := 0
	maximumTurnCount := 1000
	finalTickResult := false
	for war.IsOngoing && turnCount < maximumTurnCount {
		tickResult, err := Tick(&war, nationStatesProvider)
		assert.NoError(t, err)
		turnCount++
		finalTickResult = tickResult
	}

	assert.True(t, finalTickResult)
	assert.False(t, war.IsOngoing)
}

func TestFindOngoingWarFindsAWar(t *testing.T) {
	warAtA := NewWar("", "", "warAtA", "A")
	warAtB := NewWar("", "", "warAtB", "A")

	foundWar := FindOngoingWarAt([]War{warAtA, warAtB}, "A")

	assert.Equal(t, "warAtA", foundWar.Name)
}

func TestFindOngoingWarDoesntReturnACompletedWar(t *testing.T) {
	war := NewWar("", "", "", "A")
	war.IsOngoing = false

	foundWar := FindOngoingWarAt([]War{war}, "A")

	assert.Nil(t, foundWar)
}

func TestMorePowerfulNationDoesntAlwaysWinWar(t *testing.T) {

	nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()

	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(60)
	nationStatesProvider.PutNationData(*defender)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(40)
	nationStatesProvider.PutNationData(*attacker)

	attackerWinCount := 0
	defenderWinCount := 0
	totalLength := 0
	minimumLength := math.MaxInt32
	maximumLength := 0

	totalNumberOfSimulations := 10000

	for warIndex := 0; warIndex < totalNumberOfSimulations; warIndex++ {
		war := databasemap.NewWar(attacker.Id, defender.Id, "", "")

		length := 0
		for war.IsOngoing {
			Tick(&war, nationStatesProvider)
			length++
		}

		winnerID := WarAdvantage(war)
		if winnerID != nil && *winnerID == attacker.Id {
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
