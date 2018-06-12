package main

import (
	"fmt"
	"math/bits"
)

// Board holds pegs, indicated by their bool value
type Board uint64

// NewBoard makes a board with one peg open
func NewBoard() Board {
	var b Board
	b = ((1 << 15) - 1) ^ (1 << 4)
	return b
}

// PegCount counts pegs
func (b *Board) PegCount() int {
	return bits.OnesCount64(uint64(*b))
}

// IsComplete means the board only has one peg
func (b Board) IsComplete() bool {
	return 1 >= b.PegCount()
}

func (b Board) peg(ix uint) string {
	if (uint64(b) & (uint64(1) << ix)) > 0 {
		return "*"
	}
	return "O"
}

// Print dumps the board to the console
func (b Board) Print() {
	fmt.Printf("    %s\n", b.peg(0))
	fmt.Printf("   %s %s\n", b.peg(1), b.peg(2))
	fmt.Printf("  %s %s %s\n", b.peg(3), b.peg(4), b.peg(5))
	fmt.Printf(" %s %s %s %s\n", b.peg(6), b.peg(7), b.peg(8), b.peg(9))
	fmt.Printf("%s %s %s %s %s\n", b.peg(10), b.peg(11), b.peg(12), b.peg(13), b.peg(14))
}
