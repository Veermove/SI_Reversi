package main

import (
	"fmt"
)

func main() {
	board := InitGame()
	PrintBoard(&board)

	for _, mv := range GetLegalMoves(BLACK, &board) {
		fmt.Println(FormatMove(mv), mv.First, mv.Second)
	}

	changed := MakeMove(BLACK, &board, Pair[int, int]{5, 3})
	PrintBoard(&changed)
}
