package main

import "fmt"

// BoardMove is a tuple of a board and a move
type BoardMove struct {
	board Board
	move  Move
	prev  *BoardMove
}

func (bm *BoardMove) printMovesInner() {
	if bm.prev != nil {
		bm.prev.printMovesInner()
	}
	if bm.move.from < 255 {
		fmt.Print(bm.move)
	}
}

// PrintMoves prints the moves in the order they occurred. Not thread safe
func (bm *BoardMove) PrintMoves() {
	bm.printMovesInner()
	fmt.Println()
}

// NextStates enumerates only the valid states that can come from the given state
func (bm *BoardMove) NextStates(toCheck chan<- *BoardMove) int {
	c := 0
	bval := uint64(bm.board)
	for i := 0; i < len(moves); i++ {
		move := bitMoves[i]
		if (bval&move.from == move.from) && (bval&move.to == 0) {
			nb := Board((bval ^ move.from) | move.to)
			toCheck <- &BoardMove{nb, moves[i], bm}
			c++
		}
	}
	return c
}
