package player

import (
	"errors"
	"go-ticket-to-ride/pkg/game"
	"log/slog"
	"slices"

	"github.com/mcaci/graphgo/graph"
	"github.com/mcaci/graphgo/path"
	"github.com/mcaci/graphgo/visit"
)

type WithTickets struct {
	id         int
	ownedLines game.Board
	tickets    []game.Ticket
}

func NewWithTickets(id int, t []game.Ticket) *WithTickets {
	slog.Info("new player:", "Player", id, "tickets", t)
	return &WithTickets{id: id, tickets: t, ownedLines: graph.New[game.City](graph.ArcsListType, false)}
}

func (p *WithTickets) Play() func(game.Board) (game.City, game.City) {
	randomSelection := func(b game.Board) (game.City, game.City) {
		localBoard := graph.Copy(b)
		chosenLine, ok := PseudoRandomLine(localBoard)
		if !ok {
			return "", ""
		}
		slog.Info("pseudo-random train line choice:", "Player", p.id, "Line", chosenLine)
		chosenLine.P.(*game.TrainLineProperty).Occupy()
		p.ownedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
		p.ownedLines.AddVertex(chosenLine.X)
		p.ownedLines.AddVertex(chosenLine.Y)
		doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
			return tl.X.E == chosenLine.X.E && tl.Y.E == chosenLine.Y.E && !tl.P.(*game.TrainLineProperty).Occupied
		}, localBoard)
		if doubleLine != nil {
			doubleLine.P.(*game.TrainLineProperty).Occupy()
		}
		return chosenLine.X.E, chosenLine.Y.E
	}
	shortestPath := func(b game.Board) (game.City, game.City) {
		localBoard := graph.Copy(b)
	updatedBoard:
		for len(localBoard.Edges()) > 0 {
			ticket, err := p.NextAvailableTicket()
			if err != nil {
				return randomSelection(localBoard)
			}
			cX, cY := game.FindCity(ticket.X, localBoard), game.FindCity(ticket.Y, localBoard)
			tX, tY := (*graph.Vertex[game.City])(cX), (*graph.Vertex[game.City])(cY)
			// If there is no path between the two cities, the ticket is done and you move to the next one
			if !visit.ExistsPath(localBoard, tX, tY) {
				ticket.Done = true
				continue
			}
			// If there is a path between the two cities, you find the shortest path and take the first line
			shortest := path.Shortest(localBoard, path.BellmanFordDist(localBoard, tX), tX, tY)
			for i := 0; i < len(shortest)-1; i++ {
				chosenLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
					return tl.X.E == shortest[i].E && tl.Y.E == shortest[i+1].E || tl.X.E == shortest[i+1].E && tl.Y.E == shortest[i].E
				}, localBoard)
				chosenLineEdge := (*graph.Edge[game.City])(chosenLine)
				owned := p.ownedLines.ContainsEdge(chosenLineEdge)
				if owned {
					continue
				}
				occupiedNotOwned := chosenLine.P.(*game.TrainLineProperty).Occupied
				if occupiedNotOwned {
					localBoard.RemoveEdge(chosenLineEdge)
					continue updatedBoard
				}
				slog.Info("train line choice using shortest path:", "Player", p.id, "Line", chosenLine)
				chosenLine.P.(*game.TrainLineProperty).Occupy()
				p.ownedLines.AddVertex(chosenLine.X)
				p.ownedLines.AddVertex(chosenLine.Y)
				p.ownedLines.AddEdge((*graph.Edge[game.City])(chosenLine))
				// Remove edge for double route if needed
				doubleLine := game.FindLineFunc(func(tl *game.TrainLine) bool {
					return tl.X.E == chosenLine.X.E && tl.Y.E == chosenLine.Y.E && !tl.P.(*game.TrainLineProperty).Occupied
				}, localBoard)
				if doubleLine != nil {
					doubleLine.P.(*game.TrainLineProperty).Occupy()
				}
				// Check if ticket is done after taking the line
				if visit.ExistsPath(p.ownedLines, tX, tY) {
					ticket.Done, ticket.Ok = true, true
				}
				return chosenLine.X.E, chosenLine.Y.E
			}
			// If all the lines of the ticket are owned, the ticket is done
			ticket.Done, ticket.Ok = true, true
		}
		return "", ""
	}
	if !p.HasTicketsToComplete() {
		return randomSelection
	}
	return shortestPath
}

func (p *WithTickets) Tickets() []game.Ticket { return p.tickets }

func (p *WithTickets) TrainLines() []*graph.Edge[game.City] { return p.ownedLines.Edges() }

func (p *WithTickets) NextAvailableTicket() (*game.Ticket, error) {
	i := slices.IndexFunc(p.tickets, func(t game.Ticket) bool { return !t.Done })
	if i < 0 {
		return nil, errors.New("no available tickets")
	}
	return &p.tickets[i], nil
}

func (p *WithTickets) HasTicketsToComplete() bool {
	return slices.ContainsFunc(p.tickets, func(t game.Ticket) bool { return !t.Done })
}
