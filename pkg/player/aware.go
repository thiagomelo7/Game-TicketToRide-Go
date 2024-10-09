package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"
)

type Aware struct {
	id            int
	occupiedLines []*game.TrainLine
	tickets       []game.Ticket
}

func NewAware(id int) *Aware { return &Aware{id: id} }

func (p *Aware) Play() func(g game.Board) {
	return func(g game.Board) {
		chosenLine := g.Edges()[0]
		p.occupiedLines = append(p.occupiedLines, (*game.TrainLine)(chosenLine))
		slog.Info("New Train line taken:", "Player", p.id, "Line", *chosenLine)
		g.RemoveEdge(chosenLine)
	}
}

func (p *Aware) Score() int {
	var score int
	for i := range p.occupiedLines {
		score += p.occupiedLines[i].P.(game.TrainLineProperty).Weight()
	}
	return score
}
