package main

import (
	"fmt"

	"github.com/corentings/chess/v2"
)

func main() {
	game := chess.NewGame()

	fmt.Println("Starting position")
	fmt.Println(game.Position().Board().Draw())

	fmt.Println("FEN:", game.Position().String())
	fmt.Println("Legal moves available:", len(game.ValidMoves()))
}
