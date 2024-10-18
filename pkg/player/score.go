package player

import (
	"go-ticket-to-ride/pkg/game"

	"github.com/mcaci/graphgo/graph"
)

type Scorer interface {
	TrainLines() []*graph.Edge[game.City]
	Tickets() []game.Ticket
}

func Score(l Scorer) int {
	// ticket to ride scoring map
	scores := map[int]int{1: 1, 2: 2, 3: 4, 4: 7, 5: 10, 6: 15}
	var score int
	for _, tl := range l.TrainLines() {
		ln := tl.P.(*game.TrainLineProperty).Weight()
		score += scores[ln]
	}
	for _, t := range l.Tickets() {
		if !t.Done {
			continue
		}
		score += t.Value
	}
	return score
}
