package war

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutWarIsRetrievedFromList(t *testing.T) {
	someWar := NewWar(nil, nil, "someName", "someTerritory")
	warList := NewWarProviderSimpleList()
	warList.PutWars([]War{someWar})
	retrievedWars, err := warList.GetWars()
	assert.NoError(t, err)
	assert.Len(t, retrievedWars, 1)
	assert.Equal(t, "someName", retrievedWars[0].Name)
}

func TestUpdatedWarIsRetrievedAfterBeingUpdated(t *testing.T) {
	someWar := NewWar(nil, nil, "someName", "someTerritory")
	someWar.IsOngoing = true
	warList := NewWarProviderSimpleList()
	warList.PutWars([]War{someWar})
	someWar.IsOngoing = false
	warList.PutWars([]War{someWar})
	retrievedWars, err := warList.GetWars()
	assert.NoError(t, err)
	assert.Len(t, retrievedWars, 1)
	assert.False(t, retrievedWars[0].IsOngoing)
}
