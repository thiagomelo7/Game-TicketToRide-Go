package main

import (
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"log"
	"log/slog"

	"github.com/mcaci/graphgo/graph"
	"github.com/mcaci/graphgo/path"
)

func main() {
	b, err := data.Routes()
	if err != nil {
		log.Fatal(err)
	}
	tickets, err := data.Tickets()
	if err != nil {
		log.Fatal(err)
	}

	p1, p2 := player.NewAware(1, tickets[0]), player.NewRandom(2)
	coin := true
	// careful as some lines are double so they shouold be counted as one
	for game.FindLineFunc(func(tl *game.TrainLine) bool {
		return !tl.P.(*game.TrainLineProperty).Occupied
	}, b) != nil {
		var play func(game.Board)
		switch coin {
		case true:
			play = p1.Play()
		case false:
			play = p2.Play()
		}
		play(b)
		coin = !coin
	}
	slog.Info("end game", "Score P1", player.Score(p1), "Score P2", player.Score(p2))

	// initialChecks(b)
}

func initialChecks(b game.Board) {
	s := game.FindCity("Los Angeles", b)
	d := game.FindCity("New York", b)
	dist := path.BellmanFordDist(b, (*graph.Vertex[game.City])(s))
	log.Print(path.Shortest(b, dist, (*graph.Vertex[game.City])(s), (*graph.Vertex[game.City])(d)))
	for _, e := range b.Edges() {
		e.P.(*game.TrainLineProperty).Free()
	}
	l := game.FindLineFunc(func(tl *game.TrainLine) bool {
		return tl.X.E == "Santa Fe" && tl.Y.E == "Phoenix" || tl.X.E == "Phoenix" && tl.Y.E == "Santa Fe"
	}, b)
	l.P.(*game.TrainLineProperty).Occupy()
	nb := game.FreeRoutesBoard(b)
	ndist := path.BellmanFordDist(nb, (*graph.Vertex[game.City])(s))
	log.Print(path.Shortest(nb, ndist, (*graph.Vertex[game.City])(s), (*graph.Vertex[game.City])(d)))
}
