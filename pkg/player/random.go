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

func (p *Random) Play() func(game.Board) {
	return func(g game.Board) {
		chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
			return !tl.P.(*game.TrainLineProperty).Occupied
		}, g)
		if chosenLine == nil {
			return
		}
		p.occupiedLines = append(p.occupiedLines, (*game.TrainLine)(chosenLine))
		chosenLine.P.(*game.TrainLineProperty).Occupy()
		slog.Info("New Train line taken:", "Player", p.id, "Line", chosenLine)
		// Remove edge for double route if needed
		doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
			return tl.X.E == chosenLine.X.E &&
				tl.Y.E == chosenLine.Y.E &&
				!tl.P.(*game.TrainLineProperty).Occupied
		}, g)
		if doubleLine == nil {
			return
		}
		doubleLine.P.(*game.TrainLineProperty).Occupy()
	}
}

func (p *Random) TrainLines() []*graph.Edge[game.City] {
	lines := make([]*graph.Edge[game.City], len(p.occupiedLines))
	for i, l := range p.occupiedLines {
		lines[i] = (*graph.Edge[game.City])(l)
	}
	return lines
}
