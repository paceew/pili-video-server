package session

import (
	"fmt"
	"testing"
)

var (
	sessionID string
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestSessionServe(t *testing.T) {
	t.Run("generate session", testGenerateNewSessionID)
	t.Run("delete session", testDeleteExpiredSession)
	t.Run("session expired", testExiprideSession)
	t.Run("session reset", testReSetSession)
	t.Run("session ReExpired", testExiprideSession)
}

func testGenerateNewSessionID(t *testing.T) {
	sessionID = GenerateNewSessionId("pace")
	fmt.Printf("add session id: %v\n", sessionID)
}

func testDeleteExpiredSession(t *testing.T) {
	DeleteExpiredSession(sessionID)
}

func testExiprideSession(t *testing.T) {
	uname, ok := IsSessionExpired(sessionID)
	if ok {
		fmt.Printf("session is expired!\n")
		return
	}
	fmt.Printf("session is not expired!,uname: %v!\n", uname)
}

func testReSetSession(t *testing.T) {
	ReSetSession(sessionID, "pace")
}
