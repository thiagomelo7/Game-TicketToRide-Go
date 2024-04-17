package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mcaci/graphgo/graph"
	"github.com/mcaci/graphgo/path"
)

type Cities graph.Vertex[string]
type RailwayLines graph.Edge[string]
type RailwayLinesProperty struct {
	Cost  int
	Color color
}

type TicketToRideBoard graph.ArcsList[string]

type color int8

const (
	all color = iota
	blue
	red
	green
	yellow
	white
	pink
	orange
	black
)

type Card color

var availableTrainCars int = 40

func main() {
	b := graph.New[string](graph.ArcsListType, false)
	fillRoutes(b)
	findCity := func(name string, in graph.Graph[string]) *graph.Vertex[string] {
		for _, v := range in.Vertices() {
			if v.E != name {
				continue
			}
			return v
		}
		return nil
	}
	s := findCity("Chicago", b)
	d := findCity("Miami", b)
	dist := path.BellmanFordDist(b, s)
	log.Print(path.Shortest[string](b, dist, s, d))
	log.Print(dist[d])
}

func fillRoutes(b graph.Graph[string]) error {
	f, err := os.Open("./data/USA/routes.csv")
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	for scanner.Scan() {
		if !toRead {
			toRead = true
			continue
		}
		line := scanner.Text()
		f := strings.Split(line, ",")
		x, y, cost, col := f[0], f[1], f[2], f[3]
		X := &graph.Vertex[string]{E: x}
		Y := &graph.Vertex[string]{E: y}
		b.AddVertex(X)
		b.AddVertex(Y)
		costN, err := strconv.Atoi(cost)
		if err != nil {
			return err
		}
		mCol := map[string]color{
			"X": all,
			"B": blue,
			"R": red,
			"G": green,
			"Y": yellow,
			"W": white,
			"P": pink,
			"O": orange,
			"K": black,
		}
		P := TrainLine{col: mCol[col], cost: costN}
		b.AddEdge(&graph.Edge[string]{X: X, Y: Y, P: P})
	}
	return nil
}

type TrainLine struct {
	col  color
	cost int
}

func (t TrainLine) Weight() int { return t.cost }
