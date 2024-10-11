package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
	"github.com/mcaci/graphgo/path"
	"github.com/mcaci/graphgo/visit"
)

type Aware struct {
	id            int
	occupiedLines game.Board
	ticket        game.Ticket
}

func NewAware(id int, t game.Ticket) *Aware {
	return &Aware{id: id, ticket: t, occupiedLines: graph.New[game.City](graph.ArcsListType, false)}
}

func (p *Aware) Play() func(game.Board) {
	return func(b game.Board) {
		localBoard := graph.New[game.City](graph.ArcsListType, false)
		for _, v := range b.Vertices() {
			localBoard.AddVertex(v)
		}
		for _, e := range b.Edges() {
			localBoard.AddEdge(e)
		}
		if p.ticket.Done {
			chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
				return !tl.P.(*game.TrainLineProperty).Occupied
			}, b)
			if chosenLine == nil {
				return
			}
			p.occupiedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
			chosenLine.P.(*game.TrainLineProperty).Occupy()
			slog.Info("New Train line taken:", "Player", p.id, "Line", chosenLine)
			// Remove edge for double route if needed
			doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
				return tl.X.E == chosenLine.X.E &&
					tl.Y.E == chosenLine.Y.E &&
					!tl.P.(*game.TrainLineProperty).Occupied
			}, b)
			if doubleLine == nil {
				return
			}
			doubleLine.P.(*game.TrainLineProperty).Occupy()
			return
		}
	updatedBoard:
		for len(localBoard.Edges()) > 0 {
			tX, tY := game.FindCity(p.ticket.X, localBoard), game.FindCity(p.ticket.Y, localBoard)
			d := path.BellmanFordDist(localBoard, (*graph.Vertex[game.City])(tX))
			shortest := path.Shortest(localBoard, d, (*graph.Vertex[game.City])(tX), (*graph.Vertex[game.City])(tY))
			for i := 0; i < len(shortest)-1; i++ {
				chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
					return tl.X.E == shortest[i].E && tl.Y.E == shortest[i+1].E || tl.X.E == shortest[i+1].E && tl.Y.E == shortest[i].E

				}, localBoard)
				owned := p.occupiedLines.ContainsEdge((*graph.Edge[game.City])(chosenLine))
				if owned {
					continue
				}
				occupied := chosenLine.P.(*game.TrainLineProperty).Occupied
				if occupied {
					localBoard.RemoveEdge((*graph.Edge[game.City])(chosenLine))
					continue updatedBoard
				}
				p.occupiedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
				chosenLine.P.(*game.TrainLineProperty).Occupy()
				slog.Info("New Train line taken:", "Player", p.id, "Line", chosenLine)
				// Remove edge for double route if needed
				doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
					return tl.X.E == chosenLine.X.E &&
						tl.Y.E == chosenLine.Y.E &&
						!tl.P.(*game.TrainLineProperty).Occupied
				}, localBoard)
				if doubleLine != nil {
					doubleLine.P.(*game.TrainLineProperty).Occupy()
				}
				// Check if ticket is done after taking the line (TODO: fix)
				ep := visit.ExistsPath(p.occupiedLines, (*graph.Vertex[game.City])(tX), (*graph.Vertex[game.City])(tY))
				if ep {
					p.ticket.Done = true
				}
				return
			}
		}
	}
}

func (p *Aware) TrainLines() []*graph.Edge[game.City] { return p.occupiedLines.Edges() }
