package game

import "github.com/mcaci/graphgo/graph"

type Board graph.Graph[City]
type City string
type TrainStation graph.Vertex[City]
type TrainLine graph.Edge[City]
type TrainLineProperty struct {
	Distance int
	Color    Color
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

func (t TrainLineProperty) Weight() int { return t.Distance }

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
