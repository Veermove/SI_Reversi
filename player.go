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

func BestPlayer(maxDepth int, weights *[5]float64) moveGenerator {
	cache := make(cache)
	cache[true]  = make(map[[8][8]int]float64)
	cache[false] = make(map[[8][8]int]float64)
	if weights == nil || len(weights) == 0 {
		w := [5]float64{4, 4, 5, 5, 1}
		weights = &w
	}

	return func(board *[8][8]int, color int, turn int) Pair[int, int] {
		return FindBestMove(board, color, turn, maxDepth, cache, weights)
	}
}

func evaluateStatic(board *[8][8]int, weights *[5]float64) float64 {
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

	stonesDiff := float64(blackStones - whiteStones)
	movesDiff := float64(len(GetLegalMoves(BLACK, board)) - len(GetLegalMoves(WHITE, board)))
	edgesDiff := float64(edgesBlack - edgesWhite)
	cornersDiff := float64(cornersBlack - cornersWhite)
	keySqaresDiff := float64(keyBlack - keyWhite)

	return float64(
		movesDiff * weights[0] + edgesDiff * weights[1] + cornersDiff * weights[2] + keySqaresDiff * weights[3] + stonesDiff * weights[4],
	)
}


func minmax(cache cache, board *[8][8]int, depth int, maximazingPlayer bool, alpha float64, beta float64, weights *[5]float64) float64 {
	erl, exists := cache[maximazingPlayer][*board]

	if exists {
		return erl
	}

	if depth == 0 {
		ev := evaluateStatic(board, weights)
		cache[maximazingPlayer][*board] = ev
		return ev
	}


	if maximazingPlayer {
		value := -math.MaxFloat64
		for _, move := range GetLegalMoves(BLACK, board) {
			childBoard := MakeMove(BLACK, board, move)
			value = math.Max(value, minmax(cache, &childBoard, depth - 1, false, alpha, beta, weights))
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
			value = math.Min(value, minmax(cache, &childBoard, depth - 1, true, alpha, beta, weights))
			if value < alpha {
				break
			}
			beta = math.Min(beta, value)
		}
		return value
	}
}

func FindBestMove(board *[8][8]int, color int, turn int, depth int, cache cache, weights *[5]float64) Pair[int, int] {
	var bestMove Pair[int, int]
	if color == BLACK {
		value := -math.MaxFloat64

		for _, move := range GetLegalMoves(BLACK, board) {
			childBoard := MakeMove(BLACK, board, move)
			eval := minmax(cache, &childBoard, depth - 1, false, -math.MaxFloat64, math.MaxFloat64, weights)
			if value < eval {
				value = eval
				bestMove = move
			}
		}
	} else {
		value := math.MaxFloat64

		for _, move := range GetLegalMoves(WHITE, board) {
			childBoard := MakeMove(WHITE, board, move)
			eval := minmax(cache, &childBoard, depth - 1, true, -math.MaxFloat64, math.MaxFloat64, weights)
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

