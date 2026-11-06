package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var k float64
var F1global float64

func density(x float64) float64 {
	if x > -0.2 && x <= 1.0 {
		return 0.25 * math.Pow(x+0.2, 2.0)
	}
	if x > 1.0 && x <= 2.5 {
		return k * (x - 2.5)
	}
	return 0.0
}

func integralSecondKernel() float64 {
	v := (math.Pow(2.5-2.5, 2.0) - math.Pow(1.0-2.5, 2.0)) / 2.0
	return v
}

func numericIntegral(f func(float64) float64, a, b float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	h := (b - a) / float64(n)
	s := 0.5 * (f(a) + f(b))
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		s += f(x)
	}
	return s * h
}

func writeDat(fname string, pts [][]float64) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, p := range pts {
		fmt.Fprintf(f, "%.10f %.10f\n", p[0], p[1])
	}
	return nil
}

func writeDensity(fname string, m int, xmin, xmax float64) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	for i := 0; i <= m; i++ {
		x := xmin + (xmax-xmin)*float64(i)/float64(m)
		y := density(x)
		fmt.Fprintf(f, "%.10f %.10f\n", x, y)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	N := 2000
	nbins := 54
	xmin, xmax := -0.2, 2.5
	F1global = cdf(1.0)
	b := integralSecondKernel()
	k = (1.0 - F1global) / b
	fmt.Printf("Аналитические интегралы:\n")
	fmt.Printf(" F(1) = %.12f\n", F1global)
	fmt.Printf(" b = {1}^{2.5}∫ (x-2.5) dx = %.12f\n", b)
	fmt.Printf(" k = (1 - F1) / b = %.12f\n", k)
	F2 := k * b
	fmt.Printf(" F2 = {1}^{2.5}∫ k*(x-2.5) dx = %.12f\n", F2)
	fmt.Printf(" Fx = F1 + F2 = %.12f\n\n", F1global+F2)

	n := 20000
	F1 := numericIntegral(density, -0.2, 1.0, n)
	F2 = numericIntegral(density, 1.0, 2.5, n)
	Fx := numericIntegral(density, -0.2, 2.5, n)
	fmt.Printf("Численные интегралы:\n")
	fmt.Printf(" F(1) = {-0.2}^{1}∫ f(x) dx = %.12f\n", F1)
	fmt.Printf(" F(2.5) = {1}^{2.5}∫ f(x) dx = %.12f\n", F2)
	fmt.Printf(" F(x) = {-0.2}^{2.5}∫ f(x) dx = %.12f\n\n", Fx)

	fmt.Printf("Обратная функция: F^{-1}(0) = %.12f\n", invCdf(0.0))
	fmt.Printf("Обратная функция: F^{-1}(F1) = %.12f\n", invCdf(F1global))
	fmt.Printf("Обратная функция: F^{-1}(1) = %.12f\n\n", invCdf(1.0))

	EX_num := numericIntegral(func(x float64) float64 { return x * density(x) }, -0.2, 2.5, n)
	fmt.Printf("Численное E[X] = %.12f\n\n", EX_num)
	if err := writeDensity("density.dat", 10000, xmin, xmax); err != nil {
		panic(err)
	}
	samples := make([]float64, 0, N)
	for i := 0; i < N; i++ {
		u := rand.Float64()
		x := invCdf(u)
		samples = append(samples, x)
	}
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	var sum float64
	for _, x := range samples {
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		sum += x
	}
	mean := sum / float64(len(samples))
	fmt.Printf("Эмпирические по выборке N=%d: min=%.8f, max=%.8f, mean=%.8f\n\n", N, minX, maxX, mean)

	binWidth := (xmax - xmin) / float64(nbins)
	counts := make([]int, nbins)
	for _, x := range samples {
		if x < xmin || x > xmax {
			continue
		}
		idx := int(math.Floor((x - xmin) / binWidth))
		if idx < 0 {
			idx = 0
		}
		if idx >= nbins {
			idx = nbins - 1
		}
		counts[idx]++
	}
	histPts := make([][]float64, 0, nbins)
	var histIntegral float64 = 0.0
	for i := 0; i < nbins; i++ {
		center := xmin + (float64(i)+0.5)*binWidth
		densityVal := float64(counts[i]) / (float64(N) * binWidth)
		histPts = append(histPts, []float64{center, densityVal})
		histIntegral += densityVal * binWidth
	}
	fmt.Printf("Суммарная площадь нормированной гистограммы = %.12f\n", histIntegral)
	if err := writeDat("hist.dat", histPts); err != nil {
		panic(err)
	}
	samplePts := make([][]float64, 0, len(samples))
	for _, x := range samples {
		fx := density(x)
		if fx <= 0 {
			continue
		}
		y := rand.Float64() * fx
		samplePts = append(samplePts, []float64{x, y})
	}
	if err := writeDat("samples_on_density.dat", samplePts); err != nil {
		panic(err)
	}
}

func cdf(x float64) float64 {
	if x <= -0.2 {
		return 0.0
	}
	if x <= 1.0 {
		return math.Pow(x+0.2, 3.0) / 12.0
	}
	if x <= 2.5 {
		return F1global + (k/2.0)*(math.Pow(x-2.5, 2.0)-math.Pow(1.0-2.5, 2.0))
	}
	return 1.0
}

func invCdf(u float64) float64 {
	if u <= 0.0 {
		return -0.2
	}
	if u >= 1.0 {
		return 2.5
	}
	if u <= F1global {
		return math.Cbrt(12.0*u) - 0.2
	}
	val := 2.0*(u-F1global)/k + math.Pow(1.0-2.5, 2.0)
	if val < 0 {
		val = 0
	}
	return 2.5 - math.Sqrt(val)
}
