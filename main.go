package main

import (
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"log"
	"log/slog"
	"math/rand/v2"
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

	// For debugging purposes, we can use the same tickets
	// p1, p2 := player.NewAware(1, []game.Ticket{tickets[6]}), player.NewAware(2, []game.Ticket{tickets[16], tickets[0], tickets[25]})
	ids := rand.Perm(len(tickets))
	p1, p2 := player.NewAware(1, []game.Ticket{tickets[ids[0]]}), player.NewAware(2, []game.Ticket{tickets[ids[1]], tickets[ids[2]], tickets[ids[3]]})
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
}
