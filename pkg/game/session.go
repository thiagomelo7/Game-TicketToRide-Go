package game

import (
	"go-ticket-to-ride/pkg/utils"
	"sync"
)

type Session struct {
	ID          string
	Players     []PlayerInterface
	Board       *Board
	Turn        int
	Finished    bool
	Cards       map[Color]int
	TicketsPool *[]Ticket
}

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.RWMutex
)

func NewSession() *Session {
	return &Session{
		ID: utils.GenerateID(),
	}
}

func AddSession(s *Session) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	sessions[s.ID] = s
}

func GetSession(id string) *Session {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	return sessions[id]
}
