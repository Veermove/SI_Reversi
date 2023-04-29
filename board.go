package main

import (
	"fmt"
	"strings"
)

const (
	EMPTY = iota -1
	WHITE
	BLACK
)

func MakeMove(playColor int, board_original *[8][8]int, move Pair[int, int]) [8][8]int {
	board := *board_original
	offsets := []Pair[int, int] {
		{-1, -1},{-1,  0}, {-1,  1},
		{ 0, -1},          { 0,  1},
		{ 1, -1},{ 1,  0}, { 1,  1},
	}

	board[move.First][move.Second] = playColor

	for _, attackDir := range offsets {
		cellRow := move.First
		cellCol := move.Second

		toChange := []Pair[int, int] {}

		for {
			cellRow += attackDir.First
			cellCol += attackDir.Second

			if cellRow > 7  || cellRow < 0 || cellCol > 7 || cellCol < 0 {
				toChange = nil
				break
			}

			cell := board_original[cellRow][cellCol]

			if cell == EMPTY  {
				toChange = nil
				break
			} else if cell == playColor {
				break
			} else {
				toChange = append(toChange, Pair[int, int]{cellRow, cellCol})
			}
		}

		for _, changedCell := range toChange {
			board[changedCell.First][changedCell.Second] = playColor
		}
	}

	return board
}

func GetLegalMoves(playColor int, board *[8][8]int)  []Pair[int, int] {
	offsets := []Pair[int, int] {
		{-1, -1},{-1,  0}, {-1,  1},
		{ 0, -1},          { 0,  1},
		{ 1, -1},{ 1,  0}, { 1,  1},
	}

	moves := []Pair[int, int] {}

	for rowI := 0; rowI < 8; rowI++ {
		for colI := 0; colI < 8; colI++ {
			cellValue := board[rowI][colI]

			// This cell is occupied by other palyer or empty, so don't count moves for it
			if cellValue != playColor {
				continue
			}

			for _, attackDir := range offsets {
				neighRow := rowI + attackDir.First
				neighCol := colI + attackDir.Second

				if neighRow > 7  || neighRow < 0 || neighCol > 7 || neighCol < 0 {
					continue // neighbour out of bounds
				}

				neighbourCell := board[neighRow][neighCol]

				if neighbourCell == playColor || neighbourCell == EMPTY {
					continue // no friendly fire, no attacking air
				}

				nextFreeRow := neighRow
				nextFreeCol := neighCol
				for {
					nextFreeRow += attackDir.First
					nextFreeCol += attackDir.Second

					if nextFreeRow > 7  || nextFreeRow < 0 || nextFreeCol > 7 || nextFreeCol < 0 {
						break // out of bounds direction
					}

					targetCell := board[nextFreeRow][nextFreeCol]

					if targetCell == EMPTY {
						moves = append(moves, Pair[int, int]{nextFreeRow, nextFreeCol})
						break
					}
				}
			}

		}
	}

	return moves
}

func InitGame() [8][8]int {
	var board [8][8]int
	for r := 0; r < len(board); r++ {
		for c := 0; c < len(board[r]); c++ {
			board[r][c] = EMPTY
		}
	}
	board[Row(5)][Col('d')] = WHITE
	board[Row(5)][Col('e')] = BLACK

	board[Row(4)][Col('d')] = BLACK
	board[Row(4)][Col('e')] = WHITE

	return board
}

func PrintBoard(board *[8][8]int) {
	fmt.Print("  |")
	for i := 0; i < 8; i++ {
		fmt.Print(center(fmt.Sprintf("%c", 'a' + i), 3))
		fmt.Print("|")
	}
	fmt.Print("\n")
	fmt.Print("--+")
	fmt.Println(strings.Repeat("---+", 8))

	// for i, row := range board {
	for i := len(board) - 1; i > -1; i-- {
		row := board[i]
		fmt.Printf(" %d|", i + 1)
		for _, cell := range row {
			if cell == EMPTY {
				fmt.Print(center(" ", 3))
			} else if cell == WHITE {
				fmt.Print(" W ")
			} else {
				fmt.Print(" B ")
			}

			fmt.Print("|")
		}
		fmt.Print("\n")
		fmt.Print("--+")
		fmt.Println(strings.Repeat("---+", 8))
	}
}

func Col(column rune) int {
	return int(column) - 97
}

func Row(row int) int {
	return row - 1
}

func center(str string, width int) string {
	spaces := int(float64(width - len(str)) / 2)
	return strings.Repeat(" ", spaces) + str + strings.Repeat(" ", width - (spaces + len(str)))
}

func FormatMove(move Pair[int, int]) string {
	return fmt.Sprintf("%c%d", rune(move.Second + 97) , move.First + 1)
}

type Pair[T, U any] struct {
    First  T
    Second U
}
