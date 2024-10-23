package game

import (
	"strconv"

	"github.com/mcaci/graphgo/graph"
)

type Board graph.Graph[City]
type City string
type TrainStation graph.Vertex[City]
type TrainLine graph.Edge[City]
type TrainLineProperty struct {
	Distance int
	Color    Color
	Occupied bool
}

func FreeRoutesBoard(b Board) Board {
	frb := graph.New[City](graph.ArcsListType, false)
	for _, v := range b.Vertices() {
		frb.AddVertex(v)
	}
	for _, e := range b.Edges() {
		if e.P.(*TrainLineProperty).Occupied {
			continue
		}
		frb.AddEdge(e)
	}
	return frb
}

func FindCity(name City, in Board) *TrainStation {
	for _, v := range in.Vertices() {
		if v.E != name {
			continue
		}
		return (*TrainStation)(v)
	}
	return nil
}

func FindLineFunc(f func(*TrainLine) bool, in Board) *TrainLine {
	for _, e := range in.Edges() {
		tl := (*TrainLine)(e)
		if !f(tl) {
			continue
		}
		return tl
	}
	return nil
}

func (t TrainLine) String() string {
	return string(t.X.E) + " -> " + strconv.Itoa(t.P.(*TrainLineProperty).Distance) + " -> " + string(t.Y.E)
}

func (t TrainLineProperty) Weight() int { return t.Distance }
func (t *TrainLineProperty) Occupy()    { t.Occupied = true }
func (t *TrainLineProperty) Free()      { t.Occupied = false }

type Color int8

const (
	All Color = iota
	Blue
	Red
	Green
	Yellow
	White
	Pink
	Orange
	Black
)

type Card Color

var availableTrainCars int = 40

var totalCards = map[Color]int{
	All:    14,
	Blue:   12,
	Red:    12,
	Green:  12,
	Yellow: 12,
	White:  12,
	Pink:   12,
	Orange: 12,
	Black:  12,
}

type Ticket struct {
	X, Y  City
	Value int
	Done  bool
	Ok    bool
}

func (t Ticket) String() string {
	return string(t.X) + " -> " + string(t.Y) + " : " + strconv.Itoa(t.Value) + "."
}
