package session

import (
	"sync"
	"time"

	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
)

type Session struct {
	sessionID string
	expires   time.Time
}

type SessionManager interface {
	IsValidSession(nationName string, sessionIDString string, now time.Time) (bool, error)
	AddSession(nationName string, sessionIDString string, expires time.Time) error
	RemoveSession(nationName string) error
}

type SessionManagerSimpleMap struct {
	sessions map[string]Session
	mutex    sync.Mutex
}

func NewSessionManagerSimpleMap() SessionManagerSimpleMap {
	return SessionManagerSimpleMap{
		sessions: make(map[string]Session),
	}
}

func (manager *SessionManagerSimpleMap) IsValidSession(nationName string, sessionIDString string, now time.Time) (bool, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	now.Unix()

	foundSession, doesExist := manager.sessions[nationName]
	if !doesExist {
		return false, nil
	}

	return foundSession.sessionID == sessionIDString && foundSession.expires.After(now), nil
}

func (manager *SessionManagerSimpleMap) AddSession(nationName string, sessionIDString string, expires time.Time) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.sessions[nationName] = Session{sessionID: sessionIDString, expires: expires}

	return nil
}

func (manager *SessionManagerSimpleMap) RemoveSession(nationName string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.sessions, nationName)

	return nil
}

var simpleMapInterfaceChecker SessionManager = &SessionManagerSimpleMap{}

type SessionManagerDatabase struct {
}

func (manager *SessionManagerDatabase) IsValidSession(nationName string, sessionIDString string, now time.Time) (bool, error) {

	session, err := dynamodbwrapper.GetSession(nationName)
	if err != nil {
		return false, err
	}

	expirationDate := time.Unix(session.ExpiresAtUnixSeconds, 0)

	return session.SessionID == sessionIDString && expirationDate.After(now), nil
}

func (manager *SessionManagerDatabase) AddSession(nationName string, sessionIDString string, expires time.Time) error {

	databaseSession := dynamodbwrapper.DatabaseSession{
		NationName:           nationName,
		SessionID:            sessionIDString,
		ExpiresAtUnixSeconds: expires.Unix(),
	}

	return dynamodbwrapper.PutSession(databaseSession)
}

func (manager *SessionManagerDatabase) RemoveSession(nationName string) error {

	return dynamodbwrapper.DeleteSession(nationName)
}

var databaseInterfaceChecker SessionManager = &SessionManagerDatabase{}
