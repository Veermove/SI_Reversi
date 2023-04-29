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

type gameNode struct {
	board    [8][8]int
	children []*gameNode
}


func BestPlayer(maxDepth int) moveGenerator {
	return func(board *[8][8]int, color int) Pair[int, int] {
		return FindBestMove(board, color, maxDepth)
	}
}

func evaluateStatic(board *[8][8]int) float32 {
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

	return float32(blackStones - whiteStones)
}

func evaluate(board *[8][8]int, color int, depth int) float32 {
	if depth == 0 {
		return evaluateStatic(board)
	}

	var oppositeColor int
	if color == WHITE {
		oppositeColor = BLACK
	} else {
		oppositeColor = WHITE
	}

	afterOppositeResponse := MakeMove(oppositeColor, board, FindBestMove(board, oppositeColor, depth - 1))

	return evaluate(
		&afterOppositeResponse,
		color,
		depth - 1,
	)
}

func FindBestMove(board *[8][8]int, color int, depth int) Pair[int, int] {
	legalMoves := GetLegalMoves(color, board)

	var bestMove Pair[int, int]
	var bestScore float32

	var comparator func(sc float32, bestSc float32) bool


	if color == WHITE {
		bestScore = math.MaxFloat32
		comparator = func (sc float32, bestSc float32) bool {
			return sc < bestSc
		}
	} else {
		bestScore = -1 * math.MaxFloat32
		comparator = func (sc float32, bestSc float32) bool {
			return sc > bestSc
		}
	}

	for _, move := range legalMoves {
		afterMove := MakeMove(color, board, move)
		score := evaluate(&afterMove, color, depth)
		if comparator(score, bestScore) {
			bestScore = score
			bestMove = move
		}

	}

	return bestMove
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

