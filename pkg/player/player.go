package player

import (
	"go-ticket-to-ride/pkg/game"

	"github.com/mcaci/graphgo/graph"
)

type Player struct {
	ID         string
	Name       string
	Score      int
	Tickets    []game.Ticket
	TrainLines []*graph.Edge[game.City]
	Hand       map[game.Color]int
}

func (p *Player) GetID() string                           { return p.ID }
func (p *Player) GetName() string                         { return p.Name }
func (p *Player) GetScore() int                           { return p.Score }
func (p *Player) GetTickets() []game.Ticket               { return p.Tickets }
func (p *Player) GetTrainLines() []*graph.Edge[game.City] { return p.TrainLines }
func (p *Player) GetHand() map[game.Color]int             { return p.Hand }

var _ game.PlayerInterface = (*Player)(nil)
