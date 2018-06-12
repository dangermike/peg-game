package main

import (
	"fmt"
	"runtime"
	"sync"
)

//          0
//        1   2
//      3   4   5
//    6   7   8   9
// 10  11  12  13  14

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

var bitMoves = BitMovesFromMoves(moves)

// Empty holds nothing. This is because Golang doesn't have native sets
type Empty struct{}

var empty = Empty{}

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
				ps.NextStates(toCheck)
			}
		}()
	}

	toCheck <- &BoardMove{NewBoard(), Move{255, 255, 255}, nil}
	seenBoards := make(map[uint64]struct{})
	for bs := range toCheck {
		boardNum := uint64(bs.board)

		_, seen := seenBoards[boardNum]
		if !seen {
			pc := bs.board.PegCount()
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
