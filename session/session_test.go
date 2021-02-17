package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionManagerFindsValidSession(t *testing.T) {

	manager := NewSessionManagerSimpleMap()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	isValid, err := manager.IsValidSession("nationA", "session1", ten)
	assert.True(t, isValid)
	assert.NoError(t, err)
}

func TestExpiredSessionIsntValid(t *testing.T) {

	manager := NewSessionManagerSimpleMap()

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	manager.AddSession("nationA", "session1", ten)

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	isValid, err := manager.IsValidSession("nationA", "session1", tenTen)
	assert.False(t, isValid)
	assert.NoError(t, err)
}

func TestWrongSessionIDIsntValid(t *testing.T) {

	manager := NewSessionManagerSimpleMap()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	isValid, err := manager.IsValidSession("nationA", "session2", ten)
	assert.False(t, isValid)
	assert.NoError(t, err)
}

func TestWrongNationIDIsntValid(t *testing.T) {

	manager := NewSessionManagerSimpleMap()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	isValid, err := manager.IsValidSession("nationB", "session1", ten)
	assert.False(t, isValid)
	assert.NoError(t, err)
}

func TestSessionManagerRemovedSessionIsntValid(t *testing.T) {

	manager := NewSessionManagerSimpleMap()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)
	manager.RemoveSession("nationA")

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	isValid, err := manager.IsValidSession("nationA", "session1", ten)
	assert.False(t, isValid)
	assert.NoError(t, err)
}
