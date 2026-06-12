# Chess Arena

> ChatGPT vs. Gemini playing chess in real time — and underneath it, a Go concurrency playground built to leak, profile, and tune under load.

Chess Arena is two projects wearing one codebase. On the surface, it's a service where two large language models play chess against each other while you watch the game unfold in the browser. Underneath, it's a deliberately-built lab for going deep on Go's concurrency model: goroutines, channel orchestration, memory profiling, garbage-collection tuning, and preventing goroutine leaks under heavy load.

The trick that makes both possible is a single abstraction — a `Player` interface. A player just has to take a board position and return a move. ChatGPT and Gemini are players; so is a random-move bot, or a mock that answers instantly (and occasionally hangs on purpose). Swap what's playing, and the same engine runs either a single showcase match or thousands of concurrent games.

## How it works

### Stage 1 — The match (the product)
Two LLMs play a real game. The Go backend drives the loop: ask a player to move, validate it against the legal moves (the models *will* try illegal ones), apply it, persist it, and repeat until checkmate, stalemate, or draw. A Laravel frontend polls the backend and renders the board so you can watch the game live.

### Stage 2 — The arena (the load test)
Real API games are slow, rate-limited, and cost money — useless for stress-testing concurrency. So the arena swaps the LLMs for mock players that cost nothing and respond instantly, then runs thousands of games at once through a bounded worker pool. This is where the interesting failures live: goroutines that leak when a player hangs, memory that creeps as game histories pile up, and a garbage collector under genuine pressure. The mocks can be told to hang or error on demand, so the exact bugs worth learning to catch can be triggered deliberately and reproducibly.

## Architecture
- **Engine (Go):** the game loop, the `Player` interface, move validation, persistence, the worker pool, and the HTTP API.
- **Chess rules:** [`github.com/corentings/chess/v2`](https://github.com/corentings/chess) handles legal moves, FEN, and game-over detection — no hand-rolled chess logic.
- **Storage:** PostgreSQL. The backend writes moves as games progress; the frontend reads them back.
- **Frontend (Laravel):** polls the API and renders the board with a JS chessboard library.

## Tech stack
- Go — backend and concurrency core
- PostgreSQL — game and move storage
- Laravel + a JS chessboard library — frontend viewer
- OpenAI and Gemini APIs — the Stage 1 players

## Getting started
Requires a recent Go (1.23 or newer). Database and frontend setup arrive as those stages land — see the roadmap.

```bash
git clone https://github.com/yourname/chess-arena
cd chess-arena
go run .
```

Right now this prints the starting position and confirms the chess library is wired up correctly.

## Roadmap
- [x] Project scaffold + chess library, print a board
- [ ] Core game loop with local bots
- [ ] OpenAI and Gemini players
- [ ] Illegal-move handling and retries
- [ ] PostgreSQL persistence + HTTP API
- [ ] Laravel viewer (watch a live game)
- [ ] Arena: mock players + bounded worker pool
- [ ] Concurrency tests + goroutine-leak detection
- [ ] pprof profiling, GC tuning, and rate limiting

## Why this project
This is a learning project, so the core Go is written by hand — no AI generating the concurrency logic. AI gets used the way a tutor or reviewer would ("explain this profile," "find the leak"), but the engine itself is typed out deliberately, to keep the fundamentals sharp.

## License
MIT — see `LICENSE`.