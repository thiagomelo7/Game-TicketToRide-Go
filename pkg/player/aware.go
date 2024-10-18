package player

import (
	"go-ticket-to-ride/pkg/game"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
	"github.com/mcaci/graphgo/path"
	"github.com/mcaci/graphgo/visit"
)

type TicketAwarePlayer struct {
	id            int
	occupiedLines game.Board
	ticket        []game.Ticket
}

func NewAware(id int, t []game.Ticket) *TicketAwarePlayer {
	slog.Info("chosen tickets", "tickets", t)
	return &TicketAwarePlayer{id: id, ticket: t, occupiedLines: graph.New[game.City](graph.ArcsListType, false)}
}

func (p *TicketAwarePlayer) Play() func(game.Board) {
	rand := func(b game.Board) {
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

	graphAware := func(b game.Board) {
		localBoard := graph.Copy(b)
	updatedBoard:
		for len(localBoard.Edges()) > 0 {
			ticket := p.NextAvailableTicket()
			if ticket == nil {
				rand(localBoard)
				return
			}
			tX, tY := game.FindCity(ticket.X, localBoard), game.FindCity(ticket.Y, localBoard)
			// If there is no path between the two cities, the ticket is done and you move to the next one
			ok := visit.ExistsPath(localBoard, (*graph.Vertex[game.City])(tX), (*graph.Vertex[game.City])(tY))
			if !ok {
				ticket.Done = true
				continue
			}
			// If there is a path between the two cities, you find the shortest path and take the first line
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
				occupiedNotOwned := chosenLine.P.(*game.TrainLineProperty).Occupied
				if occupiedNotOwned {
					localBoard.RemoveEdge((*graph.Edge[game.City])(chosenLine))
					continue updatedBoard
				}
				slog.Info("New Train line chosen:", "Player", p.id, "Line", chosenLine)
				chosenLine.P.(*game.TrainLineProperty).Occupy()
				p.occupiedLines.AddVertex(chosenLine.X)
				p.occupiedLines.AddVertex(chosenLine.Y)
				p.occupiedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
				// Remove edge for double route if needed
				doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
					return tl.X.E == chosenLine.X.E && tl.Y.E == chosenLine.Y.E && !tl.P.(*game.TrainLineProperty).Occupied
				}, localBoard)
				if doubleLine != nil {
					doubleLine.P.(*game.TrainLineProperty).Occupy()
				}
				// Check if ticket is done after taking the line
				ok := visit.ExistsPath(p.occupiedLines, (*graph.Vertex[game.City])(tX), (*graph.Vertex[game.City])(tY))
				if ok {
					ticket.Done = true
					ticket.Ok = true
				}
				return
			}
		}
	}

	switch p.AllTicketsDone() {
	case false:
		return graphAware
	default:
		return rand
	}
}

func RandomLine(localBoard game.Board) (*game.TrainLine, bool) {
	chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
		return !tl.P.(*game.TrainLineProperty).Occupied
	}, localBoard)
	if chosenLine == nil {
		return nil, false
	}
	return chosenLine, true
}

func (p *TicketAwarePlayer) Tickets() []game.Ticket { return p.ticket }

func (p *TicketAwarePlayer) TrainLines() []*graph.Edge[game.City] { return p.occupiedLines.Edges() }

func (p *TicketAwarePlayer) NextAvailableTicket() *game.Ticket {
	for i, t := range p.ticket {
		if t.Done {
			continue
		}
		return &p.ticket[i]
	}
	return nil
}

func (p *TicketAwarePlayer) AllTicketsDone() bool {
	for _, t := range p.ticket {
		if t.Done {
			continue
		}
		return false
	}
	return true
}
