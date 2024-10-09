package data

import (
	"bufio"
	"go-ticket-to-ride/pkg/game"
	"os"
	"strconv"
	"strings"

	"github.com/mcaci/graphgo/graph"
)

func Routes() (game.Board, error) {
	f, err := os.Open("./pkg/data/USA/routes.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	b := graph.New[game.City](graph.AdjacencyListType, false)
	for scanner.Scan() {
		if !toRead {
			toRead = true
			continue
		}
		line := scanner.Text()
		f := strings.Split(line, ",")
		x, y, ds, cl := f[0], f[1], f[2], f[3]
		X := &graph.Vertex[game.City]{E: game.City(x)}
		Y := &graph.Vertex[game.City]{E: game.City(y)}
		b.AddVertex(X)
		b.AddVertex(Y)
		d, err := strconv.Atoi(ds)
		if err != nil {
			return nil, err
		}
		m := map[string]game.Color{
			"X": game.All,
			"B": game.Blue,
			"R": game.Red,
			"G": game.Green,
			"Y": game.Yellow,
			"W": game.White,
			"P": game.Pink,
			"O": game.Orange,
			"K": game.Black,
		}
		E := game.TrainLine{X: X, Y: Y, P: game.TrainLineProperty{Color: m[cl], Distance: d}}
		b.AddEdge((*graph.Edge[game.City])(&E))
	}
	return b, nil
}
