package strategicmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerritoryLeftAsPercentDividesAndRoundsToInteger(t *testing.T) {

	territoryA := Territory{415, 0}
	territoryB := Territory{1020, 0}
	territoryC := Territory{840, 0}

	assert.Equal(t, 27, territoryA.LeftAsPercent())
	assert.Equal(t, 66, territoryB.LeftAsPercent())
	assert.Equal(t, 55, territoryC.LeftAsPercent())
}

func TestTerritoryTopAsPercentDividesAndRoundsToInteger(t *testing.T) {

	territoryA := Territory{0, 95}
	territoryB := Territory{0, 270}
	territoryC := Territory{0, 645}

	assert.Equal(t, 13, territoryA.TopAsPercent())
	assert.Equal(t, 37, territoryB.TopAsPercent())
	assert.Equal(t, 89, territoryC.TopAsPercent())
}
