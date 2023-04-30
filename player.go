package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

type cache map[bool]map[[8][8]int]float64

func BestPlayer(maxDepth int) moveGenerator {
	cache := make(cache)
	cache[true]  = make(map[[8][8]int]float64)
	cache[false] = make(map[[8][8]int]float64)

	return func(board *[8][8]int, color int) Pair[int, int] {
		return FindBestMove(board, color, maxDepth, cache)
	}
}

func evaluateStatic(board *[8][8]int) float64 {
	whiteStones := 0
	blackStones := 0

	cornersWhite := 0
	cornersBlack := 0

	edgesWhite := 0
	edgesBlack := 0

	keyWhite := 0
	keyBlack := 0

	for rowI := 0; rowI < 8; rowI++ {
		for colI := 0; colI < 8; colI++ {
			if board[rowI][colI] == BLACK {

				//Edges control
				if rowI == 0 || rowI == 7 || colI == 7 || colI == 0 {
					edgesBlack++

					// Corner control
					if (rowI == 0 && colI == 0) || (rowI == 7 && colI == 0) || (rowI == 0 && colI == 7) || (rowI == 7 && colI == 7) {
						cornersBlack++
					}
				}


				// KeySquares
				if (rowI == 3 && colI == 3) || (rowI == 3 && colI == 4) || (rowI == 4 && colI == 3) || (rowI == 4 && colI == 4) {
					keyBlack++
				}

				blackStones++
			} else if board[rowI][colI] == WHITE {

				if rowI == 0 || rowI == 7 || colI == 7 || colI == 0 {
					edgesWhite++

					// Corner control
					if (rowI == 0 && colI == 0) || (rowI == 7 && colI == 0) || (rowI == 0 && colI == 7) || (rowI == 7 && colI == 7) {
						cornersWhite++
					}
				}

				// KeySquares
				if (rowI == 3 && colI == 3) || (rowI == 3 && colI == 4) || (rowI == 4 && colI == 3) || (rowI == 4 && colI == 4) {
					keyWhite++
				}

				whiteStones++
			}
		}
	}

	stonesDiff := blackStones - whiteStones
	movesDiff := len(GetLegalMoves(BLACK, board)) - len(GetLegalMoves(WHITE, board))
	edgesDiff := edgesBlack - edgesWhite
	cornersDiff := cornersBlack - cornersWhite
	keySqaresDiff := keyBlack - keyWhite

	return float64(
		movesDiff * 4 + edgesDiff * 4 + cornersDiff * 5 + keySqaresDiff * 5 + stonesDiff,
	)
}


func minmax(cache cache, board *[8][8]int, depth int, maximazingPlayer bool, alpha float64, beta float64) float64 {
	erl, exists := cache[maximazingPlayer][*board]

	if exists {
		return erl
	}

	if depth == 0 {
		ev := evaluateStatic(board)
		cache[maximazingPlayer][*board] = ev
		return ev
	}


	if maximazingPlayer {
		value := -math.MaxFloat64
		for _, move := range GetLegalMoves(BLACK, board) {
			childBoard := MakeMove(BLACK, board, move)
			value = math.Max(value, minmax(cache, &childBoard, depth - 1, false, alpha, beta))
			if value > beta {
				break
			}
			alpha = math.Max(alpha, value)
		}
		return value
	} else {
		value := math.MaxFloat64
		for _, move := range GetLegalMoves(WHITE, board) {
			childBoard := MakeMove(WHITE, board, move)
			value = math.Min(value, minmax(cache, &childBoard, depth - 1, true, alpha, beta))
			if value < alpha {
				break
			}
			beta = math.Min(beta, value)
		}
		return value
	}
}

func FindBestMove(board *[8][8]int, color int, depth int, cache cache) Pair[int, int] {
	var bestMove Pair[int, int]
	if color == BLACK {
		value := -math.MaxFloat64

		for _, move := range GetLegalMoves(BLACK, board) {
			childBoard := MakeMove(BLACK, board, move)
			eval := minmax(cache, &childBoard, depth - 1, false, -math.MaxFloat64, math.MaxFloat64)
			if value < eval {
				value = eval
				bestMove = move
			}
		}
	} else {
		value := math.MaxFloat64
		for _, move := range GetLegalMoves(WHITE, board) {
			childBoard := MakeMove(WHITE, board, move)
			eval := minmax(cache, &childBoard, depth - 1, true, -math.MaxFloat64, math.MaxFloat64)
			if value > eval {
				value = eval
				bestMove = move
			}
		}
	}

	return bestMove
}


func RandomPlayer(board *[8][8]int, color int) Pair[int, int] {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	legalMoves := GetLegalMoves(color, board)
	if len(legalMoves) != 0 {
		return legalMoves[r.Intn(len(legalMoves))]
	} else {
		return Pair[int, int]{}
	}

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

