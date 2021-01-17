package nationstates_api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCacheDoesntHaveNation(t *testing.T) {

	timeUntilExpiration, _ := time.ParseDuration("10m")
	cache := NewCache(timeUntilExpiration)
	someTime := time.Now()
	assert.Nil(t, cache.GetNation("someName", someTime))
}

func TestCacheWithOneRecentlyAddedNationCanFindIt(t *testing.T) {

	timeUntilExpiration, _ := time.ParseDuration("10m")
	cache := NewCache(timeUntilExpiration)
	someTime := time.Now()
	cache.AddNation("nationName", Nation{Id: "nationName"}, someTime)
	foundNation := cache.GetNation("nationName", someTime)
	assert.NotNil(t, foundNation)
	assert.Equal(t, "nationName", foundNation.Id)
}

func TestCacheWithOneExpiredNationDoesntFindIt(t *testing.T) {

	timeUntilExpiration, _ := time.ParseDuration("10m")
	cache := NewCache(timeUntilExpiration)

	tenOClock, err := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.NoError(t, err)

	cache.AddNation("nationName", Nation{Id: "nationName"}, tenOClock)

	tenThirty, err := time.Parse(time.RFC3339, "2010-10-10T10:30:00Z")
	assert.NoError(t, err)

	assert.Nil(t, cache.GetNation("nationName", tenThirty))
}

func TestCacheWithExpiredNationUpdatedWithRecentDataFindsIt(t *testing.T) {

	timeUntilExpiration, _ := time.ParseDuration("10m")
	cache := NewCache(timeUntilExpiration)

	tenOClock, err := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.NoError(t, err)

	cache.AddNation("nationName", Nation{Id: "nationName"}, tenOClock)

	tenThirty, err := time.Parse(time.RFC3339, "2010-10-10T10:30:00Z")
	assert.NoError(t, err)

	cache.AddNation("nationName", Nation{Id: "nationName"}, tenThirty)

	foundNation := cache.GetNation("nationName", tenThirty)

	assert.NotNil(t, foundNation)
	assert.Equal(t, "nationName", foundNation.Id)
}
