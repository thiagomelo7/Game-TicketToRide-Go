package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
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
		// Remove edge for double route if needed
		doubleLine := &game.TrainLine{X: chosenLine.Y, Y: chosenLine.X}
		if g.ContainsEdge((*graph.Edge[game.City])(doubleLine)) {
			g.RemoveEdge((*graph.Edge[game.City])(doubleLine))
		}
	}
}

func (p *Aware) TrainLines() []*game.TrainLine { return p.occupiedLines }
