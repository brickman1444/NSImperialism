package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionManagerFindsValidSession(t *testing.T) {

	manager := NewSessionManager()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.True(t, manager.IsValidSession("nationA", "session1", ten))
}

func TestExpiredSessionIsntValid(t *testing.T) {

	manager := NewSessionManager()

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	manager.AddSession("nationA", "session1", ten)

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.False(t, manager.IsValidSession("nationA", "session1", tenTen))
}

func TestWrongSessionIDIsntValid(t *testing.T) {

	manager := NewSessionManager()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.False(t, manager.IsValidSession("nationA", "session2", ten))
}

func TestWrongNationIDIsntValid(t *testing.T) {

	manager := NewSessionManager()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.False(t, manager.IsValidSession("nationB", "session1", ten))
}

func TestSessionManagerRemovedSessionIsntValid(t *testing.T) {

	manager := NewSessionManager()

	tenTen, _ := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	manager.AddSession("nationA", "session1", tenTen)
	manager.RemoveSession("nationA")

	ten, _ := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.False(t, manager.IsValidSession("nationA", "session1", ten))
}
