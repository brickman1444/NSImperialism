package session

import (
	"sync"
	"time"
)

type Session struct {
	sessionID string
	expires   time.Time
}

type SessionManager struct {
	sessions map[string]Session
	mutex    sync.Mutex
}

func NewSessionManager() SessionManager {
	return SessionManager{
		sessions: make(map[string]Session),
	}
}

func (manager *SessionManager) IsValidSession(nationName string, sessionIDString string, now time.Time) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	foundSession, doesExist := manager.sessions[nationName]
	if !doesExist {
		return false
	}

	return foundSession.sessionID == sessionIDString && foundSession.expires.After(now)
}

func (manager *SessionManager) AddSession(nationName string, sessionIDString string, expires time.Time) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.sessions[nationName] = Session{sessionID: sessionIDString, expires: expires}
}
