package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
)

type RandomPlayer struct {
	id            int
	occupiedLines game.Board
}

func NewRandom(id int) *RandomPlayer { return &RandomPlayer{id: id} }

func (p *RandomPlayer) Play() func(game.Board) {
	return func(b game.Board) {
		localBoard := graph.Copy(b)
		chosenLine, ok := RandomLine(localBoard)
		if !ok {
			return
		}
		slog.Info("New Train line chosen:", "Player", p.id, "Line", chosenLine)
		chosenLine.P.(*game.TrainLineProperty).Occupy()
		p.occupiedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
		p.occupiedLines.AddVertex(chosenLine.X)
		p.occupiedLines.AddVertex(chosenLine.Y)
		doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
			return tl.X.E == chosenLine.X.E && tl.Y.E == chosenLine.Y.E && !tl.P.(*game.TrainLineProperty).Occupied
		}, localBoard)
		if doubleLine != nil {
			doubleLine.P.(*game.TrainLineProperty).Occupy()
		}
	}
}
func (p *RandomPlayer) Tickets() []game.Ticket { return nil }

func (p *RandomPlayer) TrainLines() []*graph.Edge[game.City] { return p.occupiedLines.Edges() }
