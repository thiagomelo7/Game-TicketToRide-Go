package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Card color

var availableTrainCars int = 40

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

func (p *player) draw(deck map[color]int) {}
func (p *player) build() func(Board)      { return nil }

func (p *player) play() func(Board) { return nil }
