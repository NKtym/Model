package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func randomMatrix(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		for j := range m[i] {
			m[i][j] = rand.Float64() + 0.01
		}
	}
	return m
}

func normalize(m [][]float64) {
	n := len(m)
	for iter := 0; iter < 1000; iter++ {
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < n; j++ {
				sum += m[i][j]
			}
			for j := 0; j < n; j++ {
				m[i][j] /= sum
			}
		}
		for j := 0; j < n; j++ {
			sum := 0.0
			for i := 0; i < n; i++ {
				sum += m[i][j]
			}
			for i := 0; i < n; i++ {
				m[i][j] /= sum
			}
		}
	}
}

func check(m [][]float64) bool {
	n := len(m)
	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += m[i][j]
		}
		if math.Abs(sum-1) > 1e-6 {
			return false
		}
	}
	for j := 0; j < n; j++ {
		sum := 0.0
		for i := 0; i < n; i++ {
			sum += m[i][j]
		}
		if math.Abs(sum-1) > 1e-6 {
			return false
		}
	}
	return true
}

func nextState(row []float64, label string) (int, float64) {
	r := rand.Float64()
	fmt.Printf("[%s] Probabilities: ", label)
	for _, p := range row {
		fmt.Printf("%.3f ", p)
	}
	fmt.Println()
	sum := 0.0
	fmt.Printf("[%s] Cumulative:   ", label)

	for i, p := range row {
		sum += p
		fmt.Printf("%.3f ", sum)
		if r <= sum {
			fmt.Printf("\n[%s] Random=%.4f → choose state %d\n", label, r, i)
			return i, r
		}
	}
	fmt.Printf("\n[%s] Random=%.4f → choose last state\n", label, r)
	return len(row) - 1, r
}

func simulate(m [][]float64, steps int, start int, name string) ([]int, [][]int) {
	n := len(m)
	freq := make([][]int, n)
	for i := range freq {
		freq[i] = make([]int, n)
	}
	seq := make([]int, steps+1)
	seq[0] = start
	state := start
	fmt.Printf("\n--- Simulation %s ---\n", name)
	for i := 0; i < steps; i++ {
		fmt.Printf("\n[%s] Step %d: current state %d\n", name, i, state)
		next, _ := nextState(m[state], name)
		fmt.Printf("[%s] Transition: %d -> %d\n", name, state, next)
		freq[state][next]++
		state = next
		seq[i+1] = state
	}
	fmt.Printf("\n[%s] Sequence:\n", name)
	for _, v := range seq {
		fmt.Printf("%d ", v)
	}
	fmt.Println("\n")
	return seq, freq
}

func autocorr(seq []int, lag int) float64 {
	n := len(seq)
	mean := 0.0
	for _, v := range seq {
		mean += float64(v)
	}
	mean /= float64(n)
	num := 0.0
	den := 0.0
	for i := 0; i < n-lag; i++ {
		num += (float64(seq[i]) - mean) * (float64(seq[i+lag]) - mean)
	}
	for i := 0; i < n; i++ {
		den += (float64(seq[i]) - mean) * (float64(seq[i]) - mean)
	}
	return num / den
}

func saveTransitions(name string, seq []int) {
	f, _ := os.Create(name)
	defer f.Close()
	for i, v := range seq {
		fmt.Fprintf(f, "%d %d\n", i, v)
	}
}

func saveAutocorr(name string, seq []int) {
	f, _ := os.Create(name)
	defer f.Close()
	for lag := 0; lag < 25; lag++ {
		a := autocorr(seq, lag)
		fmt.Fprintf(f, "%d %f\n", lag, a)
	}
}

func printMatrix(name string, m [][]float64) {
	fmt.Println(name)
	for _, r := range m {
		for _, v := range r {
			fmt.Printf("%.4f ", v)
		}
		fmt.Println()
	}
	fmt.Println()
}

func printFreq(name string, f [][]int) {
	fmt.Println(name)
	for _, r := range f {
		fmt.Println(r)
	}
	fmt.Println()
}

func main() {
	rand.Seed(42)
	n := 20
	steps := 25
	m1 := randomMatrix(n)
	m2 := randomMatrix(n)
	normalize(m1)
	normalize(m2)
	printMatrix("Matrix A:", m1)
	printMatrix("Matrix B:", m2)
	fmt.Println("Matrix A doubly stochastic:", check(m1))
	fmt.Println("Matrix B doubly stochastic:", check(m2))
	start1 := rand.Intn(n)
	start2 := rand.Intn(n)
	fmt.Println("Start A:", start1)
	fmt.Println("Start B:", start2)
	seq1, freq1 := simulate(m1, steps, start1, "A")
	seq2, freq2 := simulate(m2, steps, start2, "B")
	printFreq("Transition frequencies A:", freq1)
	printFreq("Transition frequencies B:", freq2)
	saveTransitions("graph6_1.dat", seq1)
	saveTransitions("graph6_2.dat", seq2)
	saveAutocorr("graph7_1.dat", seq1)
	saveAutocorr("graph7_2.dat", seq2)
	fmt.Println("\nФайлы сохранены:")
	fmt.Println("graph6_1.dat, graph6_2.dat")
	fmt.Println("graph7_1.dat, graph7_2.dat")
	fmt.Println("\nПостроить графики:")
	fmt.Println("gnuplot graph6.gnu")
	fmt.Println("gnuplot graph7.gnu")
}
