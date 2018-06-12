# Peg Game
![peg game](https://images-na.ssl-images-amazon.com/images/I/418L40lLY0L.jpg)

Ever been to a Cracker Barrel? Even if you haven't you've probably seen the [peg game](https://shop.crackerbarrel.com/toys-games/games/travel-games/peg-game/606154). While I'm certain that I've beaten it, I don't remember how. Also, it was probably by accident. There are lots of [instructions](https://www.google.com/search?q=solve+triangle+peg+game) on how to solve it, none will be as fast nor as overbuild as the version here. This code will solve the puzzle in <3ms on my laptop.

## Implementation
Consider the board as a series of numbered holes that either do or do not contain a peg. This can be represented as an integer, with each bit representing a hole, 1 is full, 0 is empty. You know you are done when there is only one `1` in the number, which can be [counted efficiently](http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan). This representation does not imply any sort of spatial relationship nor jumping rules.

Valid moves are defined as a triplet of where you are jumping `from`, what you are jumping `over`, and where you are jumping `to`. For a move to be available on a given board layout, `from` and `over` must be full and `to` must be empty. While it is convenient to take these as tuples of 3 bytes, we can consider them as 2 "board" integers instead. This allows checking for validity in 2 bitwise operations, then another two to transform one board state into another.

With the ability to know we've won and the ability to know how to iterate, we have everything we need. Since this is Go, I set up a circular pipeline -- a channel for boards to iterate and a channel of boards to check. The checker is faster than the permutation step. They feed into one another, but don't grow at the same rate, which is a recipe for deadlocks. One of the channels had to be larger than the standard 1-slot to avoid deadlocks. I hate that.

While the program finishing is cool, without the list of moves we can't trust that it succeeded. The output of the program is a series of `from-over-to` peg numbers, representing the moves you have to make to be a winner like me. Boards are passed as a linked list, with the head being the most recent move made. This allows all permutations of a board to share history, avoiding unnecessary copies. There is additional code to print out the state of the board after each move (`board.Print()`), but it is currently disabled.

There are many ways to get to the same board position, no reason to check again. A set is used to hold which board positions have been visited. Doing this took the process from trying 322,324 permutations of the board to 1,651, an improvement of nearly 200x. There is a potential improvement: If we were to rotate and/or flip the board so the 3 groups of 5 pegs were in a deterministic order (e.g. fewest pegs on top, next fewest on right, most on left), we would potentially cut out additional duplicative positions. That said, at 1,651 positions it is very likely not worth the computation.

The thread/channel controls are a bit unorthodox. As this is a command-line tool we could easily just abandon the permutation threads and let them die with the parent process. The last lines of `main` are completely unnecessary cooperative shutdown. That way this code could be converted to a library for... reasons. No leaks.

## Can we go crazier?
Oh yeah.

* Because the spatial relationships are defined in code, we can do more than the triangle shape with minimal modifications
   * All of the popular peg games I found had fewer than 64 holes, so any could be modeled without changing the implementation.
   * We are not limited to real-world concepts of adjacency, nor even by only 2 dimensions. Modeling a cube would be trivial.
* 3-peg moves are a limitation of the current implementation. This need not be true.
* We are not discovering possible moves, but rather hard-coding them. This allowed me to duck spatially modeling the game board, but that's lame.
* We are doing our checks in the integer unit, which is fast. Know what would be even faster? Getting some SIMD/vector instructions going here.
* The maximum branching is the limiting factor on the channel size -- too small and we deadlock. The current value is magical. In a bad way.

## Why did you write this, weirdo?
I hadn't written any Golang in a while. Besides, I like to party. And by party I mean massively overbuild stuff in my free time. If you're nice to me maybe I'll publish the Golang version of the Chutes 'n' Ladders simulator.
