package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type moveGenerator func(*[8][8]int, int) Pair[int, int]

func StartGame(whitePlayer moveGenerator, blackPlayer moveGenerator) int {

	board := InitGame()
	turn := 0

	fmt.Println("Starting game.")
	for {
		fmt.Println("Turn:", turn)

		PrintBoard(&board)

		// black turn
		blackLegalMoves := GetLegalMoves(BLACK, &board)
		blackMove := blackPlayer(&board, BLACK)
		blackHasMove := len(blackLegalMoves) != 0 && blackMove != (Pair[int, int]{ 0, 0 })

		if blackHasMove {
			if !slices.Contains(blackLegalMoves, blackMove) {
				fmt.Println("This black move is not legal.", blackMove)
				os.Exit(1)
			}

			board = MakeMove(BLACK, &board, blackMove)
		}

		fmt.Print("\n")
		PrintBoard(&board)
		fmt.Print("\n")


		// white turn
		whiteLegalMoves := GetLegalMoves(WHITE, &board)
		whiteMove := whitePlayer(&board, WHITE)
		whiteHasMove := len(whiteLegalMoves) != 0 && whiteMove != (Pair[int, int]{ 0, 0 })


		if whiteHasMove {

			if !slices.Contains(whiteLegalMoves, whiteMove) {
				fmt.Println("This white move is not legal.", whiteMove)
				os.Exit(1)
			}

			board = MakeMove(WHITE, &board, whiteMove)
		}

		if !blackHasMove && !whiteHasMove  {
			break // game ends when no player can make legal move
		}

		turn++
	}

	whiteStones := 0
	blackStones := 0

	for rowI := 0; rowI < 8; rowI++ {
		for colI := 0; colI < 8; colI++ {
			if board[rowI][colI] == BLACK {
				blackStones++
			} else if board[rowI][colI] == WHITE {
				whiteStones++
			}
		}
	}

	fmt.Println("White stones:", whiteStones)
	fmt.Println("Black stones:", blackStones)
	if whiteStones == blackStones {
		fmt.Println("Tie!")
		return EMPTY
	} else if whiteStones > blackStones {
		fmt.Println("White wins!")
		return WHITE
	} else {
		fmt.Println("Black wins")
		return BLACK
	}
}
