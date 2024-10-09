package player

import "go-ticket-to-ride/pkg/game"

type Liner interface {
	TrainLines() []*game.TrainLine
}

func Score(l Liner) int {
	// ticket to ride scoring map
	scores := map[int]int{1: 1, 2: 2, 3: 4, 4: 7, 5: 10, 6: 15}
	var score int
	for _, tl := range l.TrainLines() {
		ln := tl.P.(game.TrainLineProperty).Weight()
		score += scores[ln]
	}
	return score
}
