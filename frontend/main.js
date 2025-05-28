import { cidades } from './cidades.js';

let gameId = null;
let state = null;

const colorNames = {
  0: "All",
  1: "Blue",
  2: "Red",
  3: "Green",
  4: "Yellow",
  5: "White",
  6: "Pink",
  7: "Orange",
  8: "Black"
};

document.getElementById('create-game').onclick = async function() {
  const res = await fetch('http://localhost:8080/api/game/create', { method: 'POST' });
  const data = await res.json();
  gameId = data.game_id;
  loadState();
};


function renderGameInfo() {
  if (!state) return;
  const info = document.getElementById('game-info');
  const player1 = state.Players.find(p => p.ID === "1");
  info.innerHTML = `
    <b>Game ID:</b> ${state.ID}<br>
    <b>Turno:</b> ${state.Turn}<br>
    <b>Status:</b> ${state.Finished ? "Finalizado" : "Em andamento"}<br>
    <b>Esse é seu objetivo:</b>
    <ul>
      ${player1.Tickets.map(t => `<li>${t.X} → ${t.Y} (${t.Value} pontos)</li>`).join('')}
    </ul>
      <b>Sua mão:</b>
  <ul>
  ${Object.entries(player1.Hand || {})
    .filter(([_, qtd]) => qtd > 0)
    .map(([cor, qtd]) => `<li>${colorNames[cor]}: ${qtd}</li>`)
    .join('')}
</ul>
  `;
}

const mapImg = document.getElementById('map-img');
const canvas = document.getElementById('map-canvas');

function resizeCanvas() {
  canvas.width = mapImg.width;
  canvas.height = mapImg.height;
}
mapImg.onload = resizeCanvas;
window.onresize = resizeCanvas;
resizeCanvas();


function drawMap() {
  resizeCanvas();
  const ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // Desenhar rotas ocupadas
  if (state && state.Board && state.Board.Edges) {
    for (const edge of state.Board.Edges) {
      if (edge.occupied) {
        console.log(edge);
        // Defina a cor da linha conforme o jogador
        let color = "gray";
        if (edge.occupied_by === "1") color = "blue";
        if (edge.occupied_by === "2") color = "green";

        // Pegue as posições das cidades
        const from = cidades[edge.from];
        const to = cidades[edge.to];
        if (from && to) {
          ctx.beginPath();
          ctx.moveTo(from.x, from.y);
          ctx.lineTo(to.x, to.y);
          ctx.strokeStyle = color;
          ctx.lineWidth = 6;
          ctx.stroke();
        }
      }
    }
  }

  // Desenhar cidades normalmente
  for (const [nome, pos] of Object.entries(cidades)) {
    ctx.beginPath();
    ctx.arc(pos.x, pos.y, 8, 0, 2 * Math.PI);
    ctx.fillStyle = "red";
    ctx.fill();
    ctx.stroke();
    ctx.fillStyle = "black";
    ctx.font = "12px Arial";
    ctx.fillText(nome, pos.x + 10, pos.y);
  }
}

let selectedCity = null;

canvas.onclick = async function(e) {
  const rect = canvas.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  for (const [nome, pos] of Object.entries(cidades)) {
    const dx = x - pos.x;
    const dy = y - pos.y;
    if (dx*dx + dy*dy < 64) { 
      if (!selectedCity) {
        selectedCity = nome;
        alert(`Cidade de origem selecionada: ${nome}. Agora clique na cidade de destino.`);
      } else if (selectedCity === nome) {
        alert("Selecione uma cidade de destino diferente da origem.");
      } else {
        try {
          await playMove(selectedCity, nome);
        } catch (err) {
          alert("Erro ao enviar jogada: " + err.message); 
        }
        selectedCity = null;
      }
      break;
    }
  }
};

async function playMove(from, to) {
      console.log("Enviando para o backend:", {
    game_id: gameId,
    player_id: "1",
    from,
    to
  });
  const res = await fetch('http://localhost:8080/api/game/play', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      game_id: gameId,
      player_id: "1",
      from,
      to
    })
  });
  if (!res.ok) {
    let msg = "Erro desconhecido";
    try {
      const data = await res.json();
      msg = data.error || JSON.stringify(data);
    } catch (e) {
      msg = await res.text();
    }
    throw new Error(msg);
  }
  loadState();
}

document.getElementById('buy-train-card').onclick = async function() {
  if (!gameId) {
    alert("Crie uma partida primeiro!");
    return;
  }
  const res = await fetch('http://localhost:8080/api/game/buy-train-card', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      game_id: gameId,
      player_id: "1"
    })
  });
  if (!res.ok) {
    let msg;
    try {
      const data = await res.json();
      msg = data.error || JSON.stringify(data);
    } catch (e) {
      msg = "Erro ao comprar carta";
    }
    alert(msg);
    return;
  }
  const data = await res.json();
  await loadState(); 
  alert("Você comprou uma carta: " + colorNames[data.drawn_color]);
};

document.getElementById('swap-ticket').onclick = async function() {
  if (!gameId) {
    alert("Crie uma partida primeiro!");
    return;
  }
  const res = await fetch('http://localhost:8080/api/game/swap-tickets', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      game_id: gameId,
      player_id: "1"
    })
  });
  if (!res.ok) {
    let msg = "Erro ao trocar ticket";
    try {
      const data = await res.json();
      msg = data.error || JSON.stringify(data);
    } catch (e) {
      msg = await res.text();
    }
    alert(msg);
    return;
  }
  await loadState();
  alert("Ticket trocado com sucesso!");
};

async function playBotTurn() {
  const action = Math.floor(Math.random() * 3);

  if (action === 0) {
    showBotMessage("Player 2 está tentando ocupar uma rota...");

    const player2 = state.Players.find(p => p.ID === "2");
    const hand = player2.Hand || {};
    const colorsWithCards = Object.entries(hand)
      .filter(([color, qtd]) => Number(qtd) > 0)
      .map(([color]) => Number(color));

    console.log("Player 2 hand:", hand);
    console.log("Player 2 colorsWithCards:", colorsWithCards);
    console.log("Edges:", state.Board && state.Board.Edges);

    let possibleMoves = [];
    if (state.Board && state.Board.Edges) {
      for (const edge of state.Board.Edges) {
        if (!edge.occupied) {
          const edgeColor = Number(edge.color);
          // Rota cinza: só pode se tiver qualquer carta colorida (exceto coringa)
          if (
            edgeColor === 0 &&
            colorsWithCards.some(cor => cor !== 0)
          ) {
            possibleMoves.push([edge.from, edge.to]);
          }
          // Rota colorida: só pode se tiver a carta da cor OU coringa (All/0)
          else if (
            edgeColor !== 0 &&
            (colorsWithCards.includes(edgeColor) || colorsWithCards.includes(0))
          ) {
            possibleMoves.push([edge.from, edge.to]);
          }
        }
      }
    }

    console.log("Possible moves found:", possibleMoves);

    if (possibleMoves.length > 0) {
      const [from, to] = possibleMoves[Math.floor(Math.random() * possibleMoves.length)];
      console.log("Bot vai ocupar a rota <>", from, to);
      try {
        await fetch('http://localhost:8080/api/game/play', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            game_id: gameId,
            player_id: "2",
            from,
            to
          })
        });
      } catch (e) {
        console.log("Error on bot play:", e);
      }
    } else {
      showBotMessage("Player 2 não pode jogar, tentando comprar carta...");
      console.log("Bot could not find a move, will try to buy a card.");
      await fetch('http://localhost:8080/api/game/buy-train-card', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          game_id: gameId,
          player_id: "2"
        })
      });
    }
  } else if (action === 1) {
    showBotMessage("Player 2 está comprando uma carta...");
    await fetch('http://localhost:8080/api/game/buy-train-card', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        game_id: gameId,
        player_id: "2"
      })
    });
  } else {
    showBotMessage("Player 2 está trocando ticket...");
    await fetch('http://localhost:8080/api/game/swap-tickets', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        game_id: gameId,
        player_id: "2"
      })
    });
  }
  await loadState();
  setTimeout(() => showBotMessage(""), 3200);
}

async function loadState() {
  if (!gameId) return;
  const res = await fetch(`http://localhost:8080/api/game/state?game_id=${gameId}`);
  state = await res.json();
  renderGameInfo();
  drawMap();

  const turno = state.Turn % state.Players.length;
  if (state.Players[turno].ID === "2" && !state.Finished) {
    setTimeout(playBotTurn, 1500);
  }
}

function showBotMessage(msg) {
  document.getElementById('bot-message').innerText = msg;
}

 loadState();