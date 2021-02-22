package war

import (
	"math"
	"testing"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/stretchr/testify/assert"
)

func TestNoOneHasAdvantageWhenWarScoreIsZero(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "", 0)

	advantage := WarAdvantage(war)

	assert.Nil(t, advantage)
}

func TestAttackerHasAdvantageWhenWarHasPositiveScore(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "", 0)
	war.Score = 1

	advantageID := WarAdvantage(war)
	assert.NotNil(t, advantageID)

	assert.Equal(t, "attacker", *advantageID)
}

func TestDefenderHasAdvantageWhenWarHasNegativeScore(t *testing.T) {
	war := databasemap.NewWar("attacker", "defender", "", "", 0)
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
	war := databasemap.NewWar("", "", "", "", 0)

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

	year := 0
	war := databasemap.NewWar(attacker.Id, defender.Id, "", "", year)

	scoreTurnZero := war.Score
	assert.Equal(t, 0, scoreTurnZero)
	year++

	Tick(&war, nationStatesProvider, year)

	scoreTurnOne := war.Score
	assert.NotEqual(t, scoreTurnZero, scoreTurnOne)
	year++

	Tick(&war, nationStatesProvider, year)

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

	year := 0
	war := databasemap.NewWar(attacker.Id, defender.Id, "", "", year)

	assert.True(t, war.IsOngoing)

	maximumTurnCount := 1000
	finalTickResult := false
	for war.IsOngoing && year < maximumTurnCount {
		tickResult, err := Tick(&war, nationStatesProvider, year)
		assert.NoError(t, err)
		year++
		finalTickResult = tickResult
	}

	assert.True(t, finalTickResult)
	assert.False(t, war.IsOngoing)
}

func TestFindOngoingWarFindsAWar(t *testing.T) {
	warAtA := databasemap.NewWar("", "", "warAtA", "A", 0)
	warAtB := databasemap.NewWar("", "", "warAtB", "A", 0)

	foundWar := FindOngoingWarAt([]databasemap.DatabaseWar{warAtA, warAtB}, "A")

	assert.Equal(t, "warAtA", foundWar.ID)
}

func TestFindOngoingWarDoesntReturnACompletedWar(t *testing.T) {
	war := databasemap.NewWar("", "", "", "A", 0)
	war.IsOngoing = false

	foundWar := FindOngoingWarAt([]databasemap.DatabaseWar{war}, "A")

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

	totalNumberOfSimulations := 10000

	for warIndex := 0; warIndex < totalNumberOfSimulations; warIndex++ {
		year := 0
		war := databasemap.NewWar(attacker.Id, defender.Id, "", "", year)

		for war.IsOngoing {
			Tick(&war, nationStatesProvider, year)
			year++
		}

		winnerID := WarAdvantage(war)
		if winnerID != nil && *winnerID == attacker.Id {
			attackerWinCount++
		} else {
			defenderWinCount++
		}
	}

	assert.Greater(t, attackerWinCount, totalNumberOfSimulations/10)
	assert.Greater(t, defenderWinCount, totalNumberOfSimulations/10)
	assert.Greater(t, attackerWinCount, defenderWinCount)
}

func TestWarsDontTakeALongTimeToResolve(t *testing.T) {

	nationStatesProvider := nationstates_api.NewNationStatesProviderSimpleMap()

	defender := &nationstates_api.Nation{Id: "Defender"}
	defender.SetDefenseForces(60)
	nationStatesProvider.PutNationData(*defender)

	attacker := &nationstates_api.Nation{Id: "Attacker"}
	attacker.SetDefenseForces(40)
	nationStatesProvider.PutNationData(*attacker)

	totalLength := 0
	minimumLength := math.MaxInt32
	maximumLength := 0

	totalNumberOfSimulations := 10000

	for warIndex := 0; warIndex < totalNumberOfSimulations; warIndex++ {
		year := 0
		war := databasemap.NewWar(attacker.Id, defender.Id, "", "", year)

		for war.IsOngoing {
			Tick(&war, nationStatesProvider, year)
			year++
		}

		totalLength = totalLength + year

		if year < minimumLength {
			minimumLength = year
		}
		if year > maximumLength {
			maximumLength = year
		}
	}

	averageLength := float32(totalLength) / float32(totalNumberOfSimulations)

	assert.Greater(t, minimumLength, 4)
	assert.Less(t, maximumLength, 25)
	assert.Less(t, averageLength, float32(9))
}
