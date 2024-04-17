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
type Lines graph.Edge[string]
type LineProperty struct {
	Cost  int
	Color color
}

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
	b, err := fillRoutes()
	if err != nil {
		log.Fatal(err)
	}
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
	tickets, err := fillTickets()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(tickets)
}

var totalCards = map[color]int{
	all:    14,
	blue:   12,
	red:    12,
	green:  12,
	yellow: 12,
	white:  12,
	pink:   12,
	orange: 12,
	black:  12,
}

type player struct {
	tickets []ticket
	cards   map[color]int
}

type ticket struct {
	x, y  string
	score int
}

func fillTickets() ([]ticket, error) {
	f, err := os.Open("./data/USA/tickets.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	var tickets []ticket
	for scanner.Scan() {
		if !toRead {
			toRead = true
			continue
		}
		line := scanner.Text()
		f := strings.Split(line, ",")
		x, y, scoreStr := f[0], f[1], f[2]
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket{x: x, y: y, score: score})
	}
	return tickets, nil
}

func fillRoutes() (graph.Graph[string], error) {
	f, err := os.Open("./data/USA/routes.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	b := graph.New[string](graph.AdjacencyListType, false)
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
			return nil, err
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
	return b, nil
}

type TrainLine struct {
	col  color
	cost int
}

func (t TrainLine) Weight() int { return t.cost }
