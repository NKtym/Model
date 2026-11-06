package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func readIntPrompt(prompt string, def int) int {
	fmt.Printf("%s (enter для %d): ", prompt, def)
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func initZeroMatrix(m, n int) [][]int {
	mat := make([][]int, m)
	for i := 0; i < m; i++ {
		mat[i] = make([]int, n)
	}
	return mat
}

func singleRunWithout(m, n, k int, rnd *rand.Rand) [][]int {
	mat := initZeroMatrix(m, n)
	if k <= 0 || n == 0 {
		return mat
	}
	perm := rnd.Perm(n)
	ptr := 0
	limitPerRow := k
	if limitPerRow > n {
		limitPerRow = n
	}
	for i := 0; i < m && ptr < n; i++ {
		for j := 0; j < limitPerRow && ptr < n; j++ {
			mat[i][perm[ptr]] = 1
			ptr++
		}
	}
	return mat
}

func singleRunWith(m, n, k int, rnd *rand.Rand) [][]int {
	mat := initZeroMatrix(m, n)
	if n == 0 || k <= 0 {
		return mat
	}
	for i := 0; i < m; i++ {
		for j := 0; j < k; j++ {
			idx := rnd.Intn(n)
			mat[i][idx] = 1
		}
	}
	return mat
}

func addMatrix(dst, src [][]int) {
	m := len(dst)
	if m == 0 {
		return
	}
	n := len(dst[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			dst[i][j] += src[i][j]
		}
	}
}

func printMatrix(mat [][]int) {
	m := len(mat)
	if m == 0 {
		fmt.Println("[]")
		return
	}
	n := len(mat[0])
	fmt.Printf("    ")
	for j := 0; j < n; j++ {
		fmt.Printf("%4d", j+1)
	}
	fmt.Println()
	for i := 0; i < m; i++ {
		fmt.Printf("%3d:", i+1)
		for j := 0; j < n; j++ {
			fmt.Printf("%4d", mat[i][j])
		}
		fmt.Println()
	}
}

func rowSums(mat [][]int) []int {
	m := len(mat)
	if m == 0 {
		return []int{}
	}
	n := len(mat[0])
	out := make([]int, m)
	for i := 0; i < m; i++ {
		s := 0
		for j := 0; j < n; j++ {
			s += mat[i][j]
		}
		out[i] = s
	}
	return out
}

func colSums(mat [][]int) []int {
	m := len(mat)
	if m == 0 {
		return []int{}
	}
	n := len(mat[0])
	out := make([]int, n)
	for j := 0; j < n; j++ {
		s := 0
		for i := 0; i < m; i++ {
			s += mat[i][j]
		}
		out[j] = s
	}
	return out
}

func main() {
	fm := flag.Int("m", 0, "")
	fn := flag.Int("n", 0, "")
	fk := flag.Int("k", 0, "")
	fr := flag.Int("r", 0, "")
	seed := flag.Int64("seed", 0, "")
	flag.Parse()

	m, n, k, R := 0, 0, 0, 0
	if *fm > 0 && *fn > 0 && *fk > 0 && *fr > 0 {
		m, n, k, R = *fm, *fn, *fk, *fr
	} else {
		defM, defN, defK, defR := 4, 6, 3, 3
		m = readIntPrompt("Введите m (число строк, 3 < m < 100)", defM)
		n = readIntPrompt("Введите n (число позиций, 3 < n < 100)", defN)
		k = readIntPrompt("Введите k (число элементов за прогон)", defK)
		R = readIntPrompt("Введите R (число прогонов)", defR)
	}

	if m <= 3 || m >= 100 || n <= 3 || n >= 100 || k <= 0 || R <= 0 {
		fmt.Println("Ошибка: проверьте параметры m, n, k, R")
		return
	}

	seedVal := *seed
	if seedVal == 0 {
		seedVal = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seedVal))
	fmt.Printf("Параметры: m=%d, n=%d, k=%d, R=%d, seed=%d\n\n", m, n, k, R, seedVal)

	fmt.Println("Массив n (элементы):")
	for i := 1; i <= n; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println("\nМассив m (строки):")
	for i := 1; i <= m; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// Режимы: без возвращения и с возвращением
	modes := []struct {
		name string
		mode string
	}{
		{"Без возвращения", "without"},
		{"С возвращением", "with"},
	}

	for _, md := range modes {
		fmt.Println("=== Режим:", md.name, "===")
		sumMat := initZeroMatrix(m, n)
		for run := 1; run <= R; run++ {
			var runMat [][]int
			if md.mode == "without" {
				runMat = singleRunWithout(m, n, k, rnd)
			} else {
				runMat = singleRunWith(m, n, k, rnd)
			}
			fmt.Printf("\nПрогон %d/%d:\n", run, R)
			printMatrix(runMat)
			addMatrix(sumMat, runMat)
		}
		fmt.Println("\nСуммарная матрица:")
		printMatrix(sumMat)
		fmt.Println("Суммы по строкам:", rowSums(sumMat))
		fmt.Println("Суммы по столбцам:", colSums(sumMat))
		fmt.Println()
	}

	// Размещение — тот же алгоритм, выводится только если m == n
	if m == n {
		fmt.Println("=== Режим: Размещение ===")
		sumMat := initZeroMatrix(m, n)
		for run := 1; run <= R; run++ {
			runMat := singleRunWithout(m, n, k, rnd)
			fmt.Printf("\nПрогон %d/%d:\n", run, R)
			printMatrix(runMat)
			addMatrix(sumMat, runMat)
		}
		fmt.Println("\nСуммарная матрица:")
		printMatrix(sumMat)
		fmt.Println("Суммы по строкам:", rowSums(sumMat))
		fmt.Println("Суммы по столбцам:", colSums(sumMat))
		fmt.Println()
	}
}
