package data

import (
	"bufio"
	"go-ticket-to-ride/pkg/game"
	"os"
	"strconv"
	"strings"
)

type MapCity map[game.City]Coord

type Coord struct {
	X int
	Y int
}

func Cities() (MapCity, error) {
	f, err := os.Open("./pkg/data/USA/cities.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var toRead bool
	m := make(MapCity)
	for scanner.Scan() {
		if !toRead {
			toRead = true
			continue
		}
		line := scanner.Text()
		f := strings.Split(line, ",")
		xS, yS, city := f[0], f[1], f[2]
		x, err := strconv.ParseFloat(xS, 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(yS, 64)
		if err != nil {
			return nil, err
		}
		m[game.City(city)] = Coord{X: int(x), Y: int(y)}
	}
	return m, nil
}
