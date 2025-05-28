package game

import "github.com/mcaci/graphgo/graph"

type PlayerInterface interface {
	GetID() string
	GetName() string
	GetScore() int
	GetTickets() []Ticket
	GetTrainLines() []*graph.Edge[City]
	GetHand() map[Color]int
}
