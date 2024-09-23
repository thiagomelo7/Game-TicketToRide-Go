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

type Board graph.Graph[City]
type City string
type TrainStation graph.Vertex[City]
type TrainLine graph.Edge[City]
type TrainLineProperty struct {
	Distance int
	Color    color
}

func (t TrainLineProperty) Weight() int { return t.Distance }

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

func main() {
	b, err := fillRoutes()
	if err != nil {
		log.Fatal(err)
	}
	findCity := func(name City, in Board) *TrainStation {
		for _, v := range in.Vertices() {
			if v.E != name {
				continue
			}
			return (*TrainStation)(v)
		}
		return nil
	}

	s := findCity("Chicago", b)
	d := findCity("Miami", b)
	dist := path.BellmanFordDist(b, (*graph.Vertex[City])(s))
	log.Print(path.Shortest[City](b, dist, (*graph.Vertex[City])(s), (*graph.Vertex[City])(d)))
	log.Print(dist[(*graph.Vertex[City])(d)])
	tickets, err := fillTickets()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(tickets)
}

func fillRoutes() (Board, error) {
	f, err := os.Open("./data/USA/routes.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	b := graph.New[City](graph.AdjacencyListType, false)
	for scanner.Scan() {
		if !toRead {
			toRead = true
			continue
		}
		line := scanner.Text()
		f := strings.Split(line, ",")
		x, y, ds, cl := f[0], f[1], f[2], f[3]
		X := &graph.Vertex[City]{E: City(x)}
		Y := &graph.Vertex[City]{E: City(y)}
		b.AddVertex(X)
		b.AddVertex(Y)
		d, err := strconv.Atoi(ds)
		if err != nil {
			return nil, err
		}
		m := map[string]color{
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
		E := TrainLine{X: X, Y: Y, P: TrainLineProperty{Color: m[cl], Distance: d}}
		b.AddEdge((*graph.Edge[City])(&E))
	}
	return b, nil
}
