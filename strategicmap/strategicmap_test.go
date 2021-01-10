package strategicmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerritoryLeftAsPercentDividesAndRoundsToInteger(t *testing.T) {

	territoryA := Territory{LeftPX: 415}
	territoryB := Territory{LeftPX: 1020}
	territoryC := Territory{LeftPX: 840}

	assert.Equal(t, 27, territoryA.LeftPercent())
	assert.Equal(t, 66, territoryB.LeftPercent())
	assert.Equal(t, 55, territoryC.LeftPercent())
}

func TestTerritoryTopAsPercentDividesAndRoundsToInteger(t *testing.T) {

	territoryA := Territory{TopPX: 95}
	territoryB := Territory{TopPX: 270}
	territoryC := Territory{TopPX: 645}

	assert.Equal(t, 13, territoryA.TopPercent())
	assert.Equal(t, 37, territoryB.TopPercent())
	assert.Equal(t, 89, territoryC.TopPercent())
}
