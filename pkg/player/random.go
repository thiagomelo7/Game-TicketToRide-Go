package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
)

type PseudoRandomPlayer struct {
	id            int
	occupiedLines game.Board
}

func NewPRPl(id int) *PseudoRandomPlayer {
	return &PseudoRandomPlayer{id: id, occupiedLines: graph.New[game.City](graph.ArcsListType, false)}
}

func (p *PseudoRandomPlayer) Play() func(game.Board) (game.City, game.City) {
	return func(b game.Board) (game.City, game.City) {
		localBoard := graph.Copy(b)
		chosenLine, ok := PseudoRandomLine(localBoard)
		if !ok {
			return "", ""
		}
		slog.Info("pseudo-random train line choice:", "Player", p.id, "Line", chosenLine)
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
		return chosenLine.X.E, chosenLine.Y.E
	}
}
func (p *PseudoRandomPlayer) Tickets() []game.Ticket { return nil }

func (p *PseudoRandomPlayer) TrainLines() []*graph.Edge[game.City] { return p.occupiedLines.Edges() }

func PseudoRandomLine(localBoard game.Board) (*game.TrainLine, bool) {
	chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
		return !tl.P.(*game.TrainLineProperty).Occupied
	}, localBoard)
	if chosenLine == nil {
		return nil, false
	}
	return chosenLine, true
}
