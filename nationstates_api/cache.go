package nationstates_api

import "time"

type CachedNation struct {
	nation         Nation
	timePulledDown time.Time
}

type Cache struct {
	internalMap         map[string]CachedNation
	timeUntilExpiration time.Duration
}

func NewCache(howLongBeforeExpires time.Duration) Cache {
	return Cache{
		internalMap:         make(map[string]CachedNation),
		timeUntilExpiration: howLongBeforeExpires,
	}
}

func (c *Cache) AddNation(nationName string, nation Nation, currentTime time.Time) {
	c.internalMap[nationName] = CachedNation{
		nation:         nation,
		timePulledDown: currentTime,
	}
}

func (c Cache) GetNation(nationName string, currentTime time.Time) *Nation {

	foundNation, didFindKey := c.internalMap[nationName]
	if !didFindKey {
		return nil
	}

	if foundNation.timePulledDown.Add(c.timeUntilExpiration).Before(currentTime) {
		return nil
	}

	return &foundNation.nation
}
