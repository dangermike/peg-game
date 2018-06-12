package main

import (
	"fmt"
	"math/bits"
	"runtime"
	"sync"
)

//          0
//        1   2
//      3   4   5
//    6   7   8   9
// 10  11  12  13  14

// Move is something that might happen
type Move struct {
	from byte
	over byte
	to   byte
}

var moves = []Move{
	Move{0, 1, 3},
	Move{1, 3, 6},
	Move{3, 6, 10},
	Move{2, 4, 7},
	Move{4, 7, 11},
	Move{5, 8, 12},
	/////
	Move{0, 2, 5},
	Move{2, 5, 9},
	Move{5, 9, 14},
	Move{1, 4, 8},
	Move{4, 8, 13},
	Move{3, 7, 12},
	/////
	Move{14, 9, 5},
	Move{9, 5, 2},
	Move{5, 2, 0},
	Move{13, 8, 4},
	Move{8, 4, 1},
	Move{12, 7, 3},
	/////
	Move{14, 13, 12},
	Move{13, 12, 11},
	Move{12, 11, 10},
	Move{9, 8, 7},
	Move{8, 7, 6},
	Move{5, 4, 3},
	/////
	Move{10, 6, 3},
	Move{6, 3, 1},
	Move{3, 1, 0},
	Move{11, 7, 4},
	Move{7, 4, 2},
	Move{12, 8, 5},
	/////
	Move{10, 11, 12},
	Move{11, 12, 13},
	Move{12, 13, 14},
	Move{6, 7, 8},
	Move{7, 8, 9},
	Move{3, 4, 5},
}

type BitMove struct {
	from uint16
	to   uint16
}

func BitMoveFromMove(move Move) BitMove {
	return BitMove{
		uint16(1<<move.from) | uint16(1<<move.over),
		1 << move.to,
	}
}

func BitMoves(moves []Move) []BitMove {
	ret := make([]BitMove, len(moves), len(moves))
	for i := 0; i < len(moves); i++ {
		ret[i] = BitMoveFromMove(moves[i])
	}
	return ret
}

var bitMoves = BitMoves(moves)

// Board holds pegs, indicated by their bool value
type Board uint16

// NewBoard makes a board with one peg open
func NewBoard() *Board {
	var b Board
	b = ((1 << 15) - 1) ^ (1 << 4)
	return &b
}

// PegCount counts pegs
func (b *Board) PegCount() int {
	return bits.OnesCount16(uint16(*b))
}

// IsComplete means the board only has one peg
func (b *Board) IsComplete() bool {
	return 1 >= b.PegCount()
}

func (b *Board) peg(ix uint) string {
	if (uint16(*b) & (uint16(1) << ix)) > 0 {
		return "*"
	}
	return "O"
}

// Print dumps the board to the console
func (b *Board) Print() {
	fmt.Printf("    %s\n", b.peg(0))
	fmt.Printf("   %s %s\n", b.peg(1), b.peg(2))
	fmt.Printf("  %s %s %s\n", b.peg(3), b.peg(4), b.peg(5))
	fmt.Printf(" %s %s %s %s\n", b.peg(6), b.peg(7), b.peg(8), b.peg(9))
	fmt.Printf("%s %s %s %s %s\n", b.peg(10), b.peg(11), b.peg(12), b.peg(13), b.peg(14))
}

// ToNumber Generates a unique number for each board
func (b *Board) ToNumber() int {
	return int(*b)
}

// BoardMove is a tuple of a board and a move
type BoardMove struct {
	board *Board
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

// Empty holds nothing. This is because Golang doesn't have native sets
type Empty struct{}

var empty = Empty{}

// NextStates enumerates only the valid states that can come from the given state
func NextStates(baseMove *BoardMove, toCheck chan<- *BoardMove) int {
	b := baseMove.board
	c := 0
	bval := uint16(*b)
	for i := 0; i < len(moves); i++ {
		move := bitMoves[i]
		if (bval&move.from == move.from) && (bval&move.to == 0) {
			nb := Board((bval ^ move.from) | move.to)
			toCheck <- &BoardMove{&nb, moves[i], baseMove}
			c++
		}
	}
	return c
}

func main() {
	workerCnt := runtime.NumCPU() - 1
	if workerCnt == 0 {
		workerCnt = 1
	}
	toCheck := make(chan *BoardMove, 1)
	toExpand := make(chan *BoardMove, 200*runtime.NumCPU())
	c := 0
	var wg sync.WaitGroup
	for x := 0; x < workerCnt; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ps := range toExpand {
				NextStates(ps, toCheck)
			}
		}()
	}

	toCheck <- &BoardMove{NewBoard(), Move{255, 255, 255}, nil}
	seenBoards := make(map[int]struct{})
	for bs := range toCheck {
		pc := bs.board.PegCount()
		boardNum := bs.board.ToNumber()

		_, seen := seenBoards[boardNum]
		// seen := false
		if !seen {
			seenBoards[boardNum] = empty
			// fmt.Printf("%d - %d: %d pegs\n", c, boardNum, pc)
			c++
			if pc > 1 {
				toExpand <- bs
			} else {
				fmt.Printf("Completed in %d moves\n", c)
				bs.PrintMoves()
				break
			}
		}
	}

	// none of this is necessary, but this is the way to gracefully shut down.
	// in case someone wants to use this as a library. hahahahaha
	close(toExpand)
	go func() {
		wg.Wait()
		close(toCheck)
	}()
	// drain toCheck while waiting for the expanders to finish
	for range toCheck {
	}
}
