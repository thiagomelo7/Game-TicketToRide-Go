package api

import (
	"encoding/json"
	"fmt"
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"math/rand"
	"net/http"

	"github.com/mcaci/graphgo/graph"
)

type EdgeDTO struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Color      int    `json:"color"`
	Occupied   bool   `json:"occupied"`
	Distance   int    `json:"distance"`
	OccupiedBy string `json:"occupied_by"`
}

func BoardToDTO(b *game.Board) []EdgeDTO {
	var edges []EdgeDTO
	for _, e := range (*b).Edges() {
		prop := e.P.(*game.TrainLineProperty)
		edges = append(edges, EdgeDTO{
			From:       string(e.X.E),
			To:         string(e.Y.E),
			Color:      int(prop.Color),
			Occupied:   prop.Occupied,
			Distance:   prop.Distance,
			OccupiedBy: prop.OccupiedBy,
		})
	}
	return edges
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	session := game.NewSession()

	tickets, err := data.Tickets()
	if err != nil {
		http.Error(w, "Erro ao carregar tickets", http.StatusInternalServerError)
		return
	}
	ticket1 := game.GetTickets(1, &tickets)
	ticket2 := game.GetTickets(1, &tickets)

	p1 := &player.Player{ID: "1", Name: "Player 1", Tickets: ticket1}
	p2 := &player.Player{ID: "2", Name: "Player 2", Tickets: ticket2}
	session.Players = []game.PlayerInterface{p1, p2}

	session.TicketsPool = &tickets

	session.Cards = make(map[game.Color]int)
	for k, v := range game.TotalCards {
		session.Cards[k] = v
	}

	board, err := data.Routes()
	if err != nil {
		http.Error(w, "Erro ao carregar rotas", http.StatusInternalServerError)
		return
	}
	session.Board = &board

	game.AddSession(session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"game_id": session.ID,
	})
}

func GetGameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	gameID := r.URL.Query().Get("game_id")
	session := game.GetSession(gameID)
	if session == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	state := map[string]interface{}{
		"ID":       session.ID,
		"Players":  session.Players,
		"Turn":     session.Turn,
		"Finished": session.Finished,
		"Board": map[string]interface{}{
			"Edges": BoardToDTO(session.Board),
		},
	}
	json.NewEncoder(w).Encode(state)
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	// Cabeçalhos essenciais para permitir CORS
	w.Header().Set("Access-Control-Allow-Origin", "*") // Permite qualquer origem
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Permite cookies/autenticação
}

func PlayMoveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		enableCORS(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}

	enableCORS(w, r)

	var req struct {
		GameID   string `json:"game_id"`
		PlayerID string `json:"player_id"`
		From     string `json:"from"`
		To       string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	session := game.GetSession(req.GameID)
	if session == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	if session.Board == nil {
		http.Error(w, "Game board not initialized", http.StatusInternalServerError)
		return
	}

	expectedPlayerID := session.Players[session.Turn%len(session.Players)].GetID()
	if req.PlayerID != expectedPlayerID {
		http.Error(w, "Not your turn", http.StatusForbidden)
		return
	}

	from := game.City(req.From)
	to := game.City(req.To)

	if game.FindCity(from, *session.Board) == nil || game.FindCity(to, *session.Board) == nil {
		http.Error(w, "Invalid city name", http.StatusBadRequest)
		return
	}

	fmt.Println("Rotas disponíveis no board:")
	for _, e := range (*session.Board).Edges() {
		tl := (*game.TrainLine)(e)
		fmt.Printf("%s <-> %s | Ocupada: %v\n", tl.X.E, tl.Y.E, tl.P.(*game.TrainLineProperty).Occupied)
	}

	line := game.FindLineFunc(func(tl *game.TrainLine) bool {
		return (tl.X.E == from && tl.Y.E == to || tl.X.E == to && tl.Y.E == from) && !tl.P.(*game.TrainLineProperty).Occupied
	}, *session.Board)
	if line == nil {
		http.Error(w, "Invalid or already occupied line", http.StatusBadRequest)
		return
	}

	corDaRota := line.P.(*game.TrainLineProperty).Color

	var playerHand map[game.Color]int
	for _, p := range session.Players {
		if p.GetID() == req.PlayerID {
			playerHand = p.GetHand()
			break
		}
	}
	if playerHand == nil {
		http.Error(w, "Player not found", http.StatusBadRequest)
		return
	}

	if corDaRota == game.All {
		// Rota cinza: pode usar qualquer carta colorida (exceto coringa), ou ALL se quiser
		found := false
		for cor, qtd := range playerHand {
			if cor != game.All && qtd > 0 {
				playerHand[cor]--
				found = true
				break
			}
		}
		// Se não achou carta colorida, tenta usar coringa (ALL)
		if !found && playerHand[game.All] > 0 {
			playerHand[game.All]--
			found = true
		}
		if !found {
			http.Error(w, "Você não tem carta para essa rota cinza", http.StatusForbidden)
			return
		}
	} else {
		// Rota colorida: pode usar carta da cor OU coringa (ALL)
		if playerHand[corDaRota] > 0 {
			playerHand[corDaRota]--
		} else if playerHand[game.All] > 0 {
			playerHand[game.All]--
		} else {
			http.Error(w, "Você não tem carta da cor necessária para essa rota", http.StatusForbidden)
			return
		}
	}

	playerHand[corDaRota]--

	line.P.(*game.TrainLineProperty).Occupy()
	line.P.(*game.TrainLineProperty).OccupiedBy = req.PlayerID

	for _, p := range session.Players {
		if p.GetID() == req.PlayerID {
			if pl, ok := p.(*player.Player); ok {
				pl.TrainLines = append(pl.TrainLines, (*graph.Edge[game.City])(line))
			}
			break
		}
	}

	session.Turn = (session.Turn + 1) % len(session.Players)

	if !game.FreeRoutesAvailable(*session.Board) {
		session.Finished = true
		fmt.Println("Fim de jogo! Não há mais rotas livres.")
		for _, p := range session.Players {
			if scorer, ok := p.(player.Scorer); ok {
				score := player.Score(scorer)
				if pl, ok := p.(*player.Player); ok {
					pl.Score = score
				}
				fmt.Printf("%s: %d pontos\n", p.GetName(), score)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func BuyTrainCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	enableCORS(w, r)

	var req struct {
		GameID   string `json:"game_id"`
		PlayerID string `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	session := game.GetSession(req.GameID)
	if session == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	expectedPlayerID := session.Players[session.Turn%len(session.Players)].GetID()
	if req.PlayerID != expectedPlayerID {
		http.Error(w, "Not your turn", http.StatusForbidden)
		return
	}

	availableColors := []game.Color{}
	for color, qty := range session.Cards {
		if qty > 0 {
			availableColors = append(availableColors, color)
		}
	}
	if len(availableColors) == 0 {
		http.Error(w, "No more cards in deck", http.StatusConflict)
		return
	}

	randomIdx := rand.Intn(len(availableColors))
	drawnColor := availableColors[randomIdx]

	session.Cards[drawnColor]--

	var hand map[game.Color]int
	for _, p := range session.Players {
		if p.GetID() == req.PlayerID {
			hand = p.GetHand()
			if hand == nil {
				hand = make(map[game.Color]int)
				if pl, ok := p.(*player.Player); ok {
					pl.Hand = hand
				}
			}
			hand[drawnColor]++
			fmt.Printf("Player %s hand after draw: %+v\n", p.GetID(), hand)
			break
		}
	}

	session.Turn = (session.Turn + 1) % len(session.Players)

	resp := map[string]interface{}{
		"drawn_color": drawnColor,
		"hand":        hand,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func SwapTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		enableCORS(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	enableCORS(w, r)

	var req struct {
		GameID   string `json:"game_id"`
		PlayerID string `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	session := game.GetSession(req.GameID)
	if session == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	if session.TicketsPool == nil || len(*session.TicketsPool) == 0 {
		http.Error(w, "No tickets available to swap", http.StatusConflict)
		return
	}

	expectedPlayerID := session.Players[session.Turn%len(session.Players)].GetID()
	if req.PlayerID != expectedPlayerID {
		http.Error(w, "Not your turn", http.StatusForbidden)
		return
	}

	for _, p := range session.Players {
		if p.GetID() == req.PlayerID {
			newTickets := game.GetTickets(1, session.TicketsPool)
			if len(newTickets) == 0 {
				http.Error(w, "No tickets available to swap", http.StatusConflict)
				return
			}
			if pl, ok := p.(*player.Player); ok {
				pl.Tickets = newTickets
			}
			break
		}
	}

	session.Turn = (session.Turn + 1) % len(session.Players)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
