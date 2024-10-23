package main

import (
	"flag"
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"

	"github.com/StephaneBunel/bresenham"
)

func main() {
	show := flag.Bool("show", false, "show the game on the map")
	flag.Parse()

	routes, err := data.Routes()
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

	var layer *image.NRGBA
	if *show {
		sourceImg, err := os.Open("./pkg/data/USA/USA_map.jpg")
		if err != nil {
			slog.Error("error occurred", "err", err)
		}
		src, err := jpeg.Decode(sourceImg)
		if err != nil {
			slog.Error("error occurred", "err", err)
		}
		layer = image.NewNRGBA(src.Bounds())
		draw.Draw(layer, layer.Bounds(), src, image.Point{}, draw.Over)
	}

	// For debugging purposes, we can use the fix the tickets
	// p1, p2 := player.NewTAPl(1, []game.Ticket{tickets[3], tickets[26], tickets[21]}), player.NewTAPl(2, []game.Ticket{tickets[12], tickets[2], tickets[22]})
	ids := rand.Perm(len(tickets))
	p1, p2 := player.NewTAPl(1, []game.Ticket{tickets[ids[0]], tickets[ids[2]], tickets[ids[4]]}), player.NewTAPl(2, []game.Ticket{tickets[ids[1]], tickets[ids[3]], tickets[ids[5]]})
	var coin bool
	var frames []*image.Paletted
	// careful as some lines are double so they shouold be counted as one
	for game.FindLineFunc(func(tl *game.TrainLine) bool {
		return !tl.P.(*game.TrainLineProperty).Occupied
	}, routes) != nil {
		var play func(game.Board) (game.City, game.City)
		coin = !coin
		switch coin {
		case true:
			play = p1.Play()
		case false:
			play = p2.Play()
		}
		a, b := play(routes)

		if *show {
			var c color.Color
			switch coin {
			case true:
				c = color.RGBA{R: 0, G: 0, B: 255, A: 255}
			case false:
				c = color.RGBA{R: 255, G: 0, B: 255, A: 255}
			}
			bresenham.DrawLine(layer, cities[a].X, cities[a].Y, cities[b].X, cities[b].Y, c)
			p := image.NewPaletted(layer.Bounds(), palette.Plan9)
			draw.Draw(p, p.Bounds(), layer, image.Point{}, draw.Over)
			frames = append(frames, p)
		}
	}

	slog.Info("end game", "Score P1", player.Score(p1), "Score P2", player.Score(p2))
	if *show {
		out, err := os.Create("./pkg/data/USA/USA_map_out.jpg")
		if err != nil {
			slog.Error("error occurred", "err", err)
		}
		jpeg.Encode(out, layer, nil)
		g := gif.GIF{
			Image: frames,
			Delay: slices.Repeat([]int{45}, len(frames)),
		}
		outGif, err := os.Create("./pkg/data/USA/USA_map_out.gif")
		if err != nil {
			slog.Error("error occurred", "err", err)
		}
		gif.EncodeAll(outGif, &g)
	}
}
