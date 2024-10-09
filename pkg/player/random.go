package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
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
		// Remove edge for double route if needed
		doubleLine := &game.TrainLine{X: chosenLine.Y, Y: chosenLine.X}
		if g.ContainsEdge((*graph.Edge[game.City])(doubleLine)) {
			g.RemoveEdge((*graph.Edge[game.City])(doubleLine))
		}
	}
}

func (p *Random) TrainLines() []*game.TrainLine { return p.occupiedLines }
