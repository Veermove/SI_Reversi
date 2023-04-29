package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

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

		PrintBoard(&board)

		// white turn
		whiteLegalMoves := GetLegalMoves(WHITE, &board)

		if len(whiteLegalMoves) != 0 {
			whiteMove := whitePlayer(&board, WHITE)

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


func RandomPlayer(board *[8][8]int, color int) Pair[int, int] {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	legalMoves := GetLegalMoves(color, board)
	return legalMoves[r.Intn(len(legalMoves))]
}

func StdinPlayer(board *[8][8]int, color int) Pair[int, int] {
	reader := bufio.NewReader(os.Stdin)
	legalMoves := GetLegalMoves(color, board)
	if color == WHITE {
		fmt.Println("White to play")
	} else {
		fmt.Println("Black to play")
	}

	for {
		text, _ := reader.ReadString('\n')
		if len(text) != 3 {
			fmt.Println("Len Error", len(text))
			continue
		}
		row, err := strconv.Atoi(string(text[1]))
		if err != nil {
			fmt.Println("Error")
			continue
		}
		move := Pair[int, int]{ Row(row), Col(rune(text[0])) }
		if slices.Contains(legalMoves, move) {
			return move
		} else {
			fmt.Println("Illegal move!")
		}
	}
}
