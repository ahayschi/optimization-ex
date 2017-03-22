package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dustin/go-humanize"
)

// Board contains a 2d array of intergers representing its state
type Board struct {
	size  int
	state [][]int
}

// NewBoard initializes a new Board of size n*n with pieces set
// to [0:n*n] in ascending order
func NewBoard(n int) *Board {
	state := make([][]int, n)
	for i := 0; i < n; i++ {
		state[i] = make([]int, n)
		for j := 0; j < n; j++ {
			state[i][j] = (i * n) + j
		}
	}

	return &Board{size: n, state: state}
}

// NewBoardRand initializes a new Board of size n*n with pieces set
// to [0:n*n] in random order
func NewBoardRand(n int) *Board {
	state := make([][]int, n)
	pieces := rand.Perm(n * n)

	for i := 0; i < n; i++ {
		state[i] = pieces[i*n : (i*n)+n]
	}

	return &Board{size: n, state: state}
}

// Copy returns a copy of the board
func (b *Board) Copy() *Board {
	s := b.size
	state := make([][]int, s)
	for i := 0; i < s; i++ {
		state[i] = make([]int, s)
		copy(state[i], b.state[i])
	}

	return &Board{size: s, state: state}
}

func (b *Board) String() string {
	var s string
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			s += fmt.Sprintf("%2d ", b.state[i][j])
		}
		s += "\n"
	}

	return s
}

// DiffN returns the number of pieces in b that are off their
// target tiles (f(x) -> h)
func (b *Board) DiffN(target *Board) int {
	if b.size != target.size {
		panic("boards must be of the same size")
	}

	c := 0
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.state[i][j] == 0 {
				continue
			}

			if b.state[i][j] != target.state[i][j] {
				c++
			}
		}
	}

	return c
}

func (b *Board) tileIndx(n int) (int, int, bool) {
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.state[i][j] == n {
				return i, j, true
			}
		}
	}

	return -1, -1, false
}

// FindNeighbors returns a slice of boards with switched
// states on tile with the x value.
func (b *Board) FindNeighbors(x int) []*Board {
	i, j, ok := b.tileIndx(x)
	if !ok {
		return nil
	}

	var ns []*Board

	// Up
	if i-1 >= 0 {
		c := b.Copy()
		c.state[i][j], c.state[i-1][j] = b.state[i-1][j], b.state[i][j]
		ns = append(ns, c)
	}

	// Down
	if i+1 < b.size {
		c := b.Copy()
		c.state[i][j], c.state[i+1][j] = b.state[i+1][j], b.state[i][j]
		ns = append(ns, c)
	}

	// Left
	if j-1 >= 0 {
		c := b.Copy()
		c.state[i][j], c.state[i][j-1] = b.state[i][j-1], b.state[i][j]
		ns = append(ns, c)
	}

	// Right
	if j+1 < b.size {
		c := b.Copy()
		c.state[i][j], c.state[i][j+1] = b.state[i][j+1], b.state[i][j]
		ns = append(ns, c)
	}

	return ns
}

// FindNeighborRand returns a random neighbor of tile == x
func (b *Board) FindNeighborRand(x int) *Board {
	ns := b.FindNeighbors(x)
	if ns == nil {
		return nil
	}

	i := rand.Intn(len(ns))
	return ns[i]
}

// MinDiffHC uses a hill-climbing algorithm to find the minimum*
// number of moves to go from b to target.
//
// *Not guaranteed to be the global minimum, most often a local
func (b *Board) MinDiffHC(target *Board, ch chan int) {
	currentBoard := b

	for {
		var nextBoard *Board
		h := math.MaxInt32 // minimize h

		candidates := currentBoard.FindNeighbors(0)
		for _, c := range candidates {
			diff := c.DiffN(target)
			if diff < h {
				nextBoard = c
				h = diff
			}
		}

		currentH := currentBoard.DiffN(target)
		if h >= currentH {
			ch <- currentH
			return
		}

		currentBoard = nextBoard
	}
}

// SimulatedAnnealParams are the tuning parameters for
// the simulated annealing experiment.
type SimulatedAnnealParams struct {
	Objective      *Board
	TemperatureMin float64
	Alpha          float64
	Iterations     int
	MaxTime        time.Duration
}

// AcceptCandidate is a standard probability function for evaulating
// candidate states for simulated annealing.
func AcceptCandidate(current, candidate int, temp float64) bool {
	// Metropolis criterion
	if candidate < current {
		return true
	}

	if temp == 0 {
		return false
	}

	d := float64(candidate - current)
	p := math.Exp(-d / temp)
	return rand.Float64() < p
}

// MinDiffSA uses a simulated annealing algorithm to find the minimum*
// number of moves to go from b to target.
//
// *Could be global minimum or a close approximation
func (b *Board) MinDiffSA(p *SimulatedAnnealParams, ch chan int) {
	currentBoard := b
	currentDiff := b.DiffN(p.Objective)
	t := 1.00
	tMin := p.TemperatureMin
	alpha := p.Alpha
	minDiff := currentDiff

	start := time.Now()
	for t > tMin {
		for i := 0; i < p.Iterations; i++ {
			candidateBoard := currentBoard.FindNeighborRand(0)
			candidateDiff := candidateBoard.DiffN(p.Objective)
			if candidateDiff < minDiff {
				minDiff = candidateDiff
			}

			if AcceptCandidate(currentDiff, candidateDiff, t) {
				currentBoard = candidateBoard
				currentDiff = candidateDiff
			}
		}
		t *= alpha

		if time.Since(start) > p.MaxTime {
			if minDiff < currentDiff {
				ch <- minDiff
			}
			ch <- currentDiff
			return
		}
	}

	if minDiff < currentDiff {
		ch <- minDiff
	}
	ch <- currentDiff
	return
}

// RunHillClimb runs an example experiment for hill-climbing
func RunHillClimb() {
	size := 3
	n := 5000000
	b1 := NewBoard(size)
	ch := make(chan int, n)
	diffs := make(map[int]int, size)
	for i := 0; i < size*size; i++ {
		diffs[i] = 0
	}

	start := time.Now()
	for i := 0; i < n; i++ {
		b2 := NewBoardRand(size)
		go b2.MinDiffHC(b1, ch)
	}

	for i := 0; i < n; i++ {
		m := <-ch
		diffs[m]++
	}

	runtime := time.Since(start).Seconds()
	sn := humanize.Comma(int64(n))
	fmt.Printf("Board Size: %d\nIterations: %s\nRuntime: %.3fs\n", size, sn, runtime)
	fmt.Println("Min\tPercent\tCount")
	for i := size*size - 1; i >= 0; i-- {
		d := diffs[i]
		per := (float64(d) / float64(n)) * 100
		sd := humanize.Comma(int64(d))
		fmt.Printf("%d\t%.3f\t%s\n", i, per, sd)
	}
}

// RunSimAnneal runs an example experiment for simulated annealing
func RunSimAnneal() {
	size := 3
	n := 100
	maxTime := 10 * time.Second
	b1 := NewBoard(size)
	ch := make(chan int, n)
	diffs := make(map[int]int, size)
	for i := 0; i < size*size; i++ {
		diffs[i] = 0
	}

	params := SimulatedAnnealParams{
		Objective:      b1,
		TemperatureMin: 0.00000000001,
		Alpha:          0.99,
		Iterations:     1000,
		MaxTime:        maxTime,
	}

	start := time.Now()
	for i := 0; i < n; i++ {
		b2 := NewBoardRand(size)
		go b2.MinDiffSA(&params, ch)
	}

	for i := 0; i < n; i++ {
		m := <-ch
		diffs[m]++
	}

	runtime := time.Since(start).Seconds()
	sn := humanize.Comma(int64(n))
	fmt.Printf("Board Size: %d\nIterations: %s\nRuntime: %.3fs\n", size, sn, runtime)
	fmt.Println("Min\tPercent\tCount")
	for i := size*size - 1; i >= 0; i-- {
		d := diffs[i]
		per := (float64(d) / float64(n)) * 100
		sd := humanize.Comma(int64(d))
		fmt.Printf("%d\t%.3f\t%s\n", i, per, sd)
	}

}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	fmt.Println("Running hill-climb and simulated annealing optimizations...")
	fmt.Printf("Running hill-climb...\n")
	RunHillClimb()

	fmt.Printf("\nRunning simulated annealing...\n")
	RunSimAnneal()
}
