package data

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ticket struct {
	X, Y  string
	Score int
}

func Tickets() ([]ticket, error) {
	f, err := os.Open("./pkg/data/USA/tickets.csv")
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
		tickets = append(tickets, ticket{X: x, Y: y, Score: score})
	}
	return tickets, nil
}
