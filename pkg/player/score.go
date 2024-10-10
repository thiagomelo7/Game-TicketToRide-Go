package player

import (
	"go-ticket-to-ride/pkg/game"

	"github.com/mcaci/graphgo/graph"
)

type Liner interface {
	TrainLines() []*graph.Edge[game.City]
}

func Score(l Liner) int {
	// ticket to ride scoring map
	scores := map[int]int{1: 1, 2: 2, 3: 4, 4: 7, 5: 10, 6: 15}
	var score int
	for _, tl := range l.TrainLines() {
		ln := tl.P.(*game.TrainLineProperty).Weight()
		score += scores[ln]
	}
	return score
}
