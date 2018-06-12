package main

// Move is something that might happen
type Move struct {
	from byte
	over byte
	to   byte
}

// BitMove is a bitwise representation of a move
type BitMove struct {
	from uint64
	to   uint64
}

// BitMoveFromMove converts a Move to a BitMove
func BitMoveFromMove(move Move) BitMove {
	return BitMove{
		uint64(1<<move.from) | uint64(1<<move.over),
		1 << move.to,
	}
}

// BitMovesFromMoves converts a slice of Move to a slice of BitMove
func BitMovesFromMoves(moves []Move) []BitMove {
	ret := make([]BitMove, len(moves), len(moves))
	for i := 0; i < len(moves); i++ {
		ret[i] = BitMoveFromMove(moves[i])
	}
	return ret
}
