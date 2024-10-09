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
	initialChecks(b)

	p1, p2 := player.NewRandom(1), player.NewRandom(2)
	var coin bool
	// careful as some lines are double so they shouold be counted as one
	for len(b.Edges()) > 0 {
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
	slog.Info("end game", "Score P1", p1.Score(), "Score P2", p2.Score())
}

func initialChecks(b game.Board) {
	s := game.FindCity("Chicago", b)
	d := game.FindCity("Miami", b)
	dist := path.BellmanFordDist(b, (*graph.Vertex[game.City])(s))
	log.Print(path.Shortest[game.City](b, dist, (*graph.Vertex[game.City])(s), (*graph.Vertex[game.City])(d)))
	log.Print(dist[(*graph.Vertex[game.City])(d)])
	tickets, err := data.Tickets()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(tickets)
}
