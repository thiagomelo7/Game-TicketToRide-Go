package main

import (
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log/slog"
	"math/rand/v2"
	"os"

	"github.com/StephaneBunel/bresenham"
)

func main() {
	b, err := data.Routes()
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	tickets, err := data.Tickets()
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	cities, err := data.Cities()
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	sourceImg, err := os.Open("./pkg/data/USA/USA_map.jpg")
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	src, err := jpeg.Decode(sourceImg)
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	out, err := os.Create("./pkg/data/USA/USA_map_out.jpg")
	if err != nil {
		slog.Error("error occurred", "err", err)
	}
	layer := image.NewNRGBA(src.Bounds())
	draw.Draw(layer, layer.Bounds(), src, image.Point{}, draw.Over)

	// For debugging purposes, we can use the same tickets
	// p1, p2 := player.NewAware(1, []game.Ticket{tickets[6]}), player.NewAware(2, []game.Ticket{tickets[16], tickets[0], tickets[25]})
	ids := rand.Perm(len(tickets))
	p1, p2 := player.NewAware(1, []game.Ticket{tickets[ids[0]]}), player.NewAware(2, []game.Ticket{tickets[ids[1]], tickets[ids[2]], tickets[ids[3]]})
	coin := true
	// careful as some lines are double so they shouold be counted as one
	for game.FindLineFunc(func(tl *game.TrainLine) bool {
		return !tl.P.(*game.TrainLineProperty).Occupied
	}, b) != nil {
		var play func(game.Board) (game.City, game.City)
		var c color.Color
		switch coin {
		case true:
			play = p1.Play()
			c = color.RGBA{0, 0, 0, 255}
		case false:
			play = p2.Play()
			c = color.RGBA{255, 255, 255, 255}
		}
		a, b := play(b)
		bresenham.DrawLine(layer, cities[a].X, cities[a].Y, cities[b].X, cities[b].Y, c)
		coin = !coin
	}
	slog.Info("end game", "Score P1", player.Score(p1), "Score P2", player.Score(p2))
	jpeg.Encode(out, layer, nil)

}
