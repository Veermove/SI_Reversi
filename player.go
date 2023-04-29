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


func BestPlayer(maxDepth int) moveGenerator {
	return func(board *[8][8]int, color int) Pair[int, int] {
		legalMoves := GetLegalMoves(color, board)
		bestMove := FindBestMove(legalMoves, board, color)
		return bestMove
	}
}

func FindBestMove(legalMoves []Pair[int, int], board *[8][8]int, color int) Pair[int, int] {
	panic("unimplemented")
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

