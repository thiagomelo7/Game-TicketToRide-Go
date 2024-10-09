package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"
)

type Random struct {
	id            int
	occupiedLines []*game.TrainLine
}

func NewRandom(id int) *Random { return &Random{id: id} }

func (p *Random) Play() func(g game.Board) {
	return func(g game.Board) {
		chosenLine := g.Edges()[0]
		p.occupiedLines = append(p.occupiedLines, (*game.TrainLine)(chosenLine))
		slog.Info("New Train line taken:", "Player", p.id, "Line", *chosenLine)
		g.RemoveEdge(chosenLine)
	}
}

func (p *Random) Score() int {
	var score int
	for i := range p.occupiedLines {
		score += p.occupiedLines[i].P.(game.TrainLineProperty).Weight()
	}
	return score
}
