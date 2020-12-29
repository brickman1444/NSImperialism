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

	war := &War{defender, attacker, 0}

	assert.Equal(t, 50, war.ScoreChangePerYear())
}
