package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type moveGenerator func(*[8][8]int, int) Pair[int, int]

func StartGame(whitePlayer moveGenerator, blackPlayer moveGenerator) {

	board := InitGame()
	turn := 0

	for {
		// black turny
		blackLegalMoves := GetLegalMoves(BLACK, &board)

		if len(blackLegalMoves) != 0 {
			blackMove := blackPlayer(&board, BLACK)

			if !slices.Contains(blackLegalMoves, blackMove) {
				fmt.Println("This move is not legal.")
				os.Exit(1)
			}

			board = MakeMove(BLACK, &board, blackMove)
		}

		// white turn
		whiteLegalMoves := GetLegalMoves(WHITE, &board)

		if len(whiteLegalMoves) != 0 {
			whiteMove := blackPlayer(&board, WHITE)

			if !slices.Contains(whiteLegalMoves, whiteMove) {
				fmt.Println("This move is not legal.")
				os.Exit(1)
			}

			board = MakeMove(WHITE, &board, whiteMove)
		}

		if len(whiteLegalMoves) == 0 && len(blackLegalMoves) == 0 {
			break // game ends when no player can make legal move
		}

		turn++
	}
}
