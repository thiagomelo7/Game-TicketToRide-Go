package data

import (
	"bufio"
	"go-ticket-to-ride/pkg/game"
	"os"
	"strconv"
	"strings"
)

func Tickets() ([]game.Ticket, error) {
	f, err := os.Open("./pkg/data/USA/tickets.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	var tickets []game.Ticket
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
		tickets = append(tickets, game.Ticket{X: game.City(x), Y: game.City(y), Value: score})
	}
	return tickets, nil
}
