package main

import "testing"

func buildTestGraph(n int, isolated int) ([]Point, []Edge, Constraints) {

	points := make([]Point, n)

	for i := 0; i < n; i++ {
		points[i] = Point{
			X: float64(i * 3),
			Y: float64(i * 5),
		}
	}

	constraints := Constraints{
		UseMaxDeg:    true,
		MaxDegree:    6,
		UseDistRange: true,
		DMin:         0,
		DMax:         100,
		NumIsolated:  isolated,
	}

	edges := buildGraph(points, 0.1, 2.0, 1, constraints)

	return points, edges, constraints
}

func TestNoSelfLoops(t *testing.T) {

	_, edges, _ := buildTestGraph(100, 0)

	for _, e := range edges {
		if e.U == e.V {
			t.Fatalf("self-loop detected: %+v", e)
		}
	}
}

func TestNoParallelEdges(t *testing.T) {

	_, edges, _ := buildTestGraph(100, 0)

	seen := map[[2]int]bool{}

	for _, e := range edges {

		u := e.U
		v := e.V

		if u > v {
			u, v = v, u
		}

		key := [2]int{u, v}

		if seen[key] {
			t.Fatalf("parallel edge detected: %+v", e)
		}

		seen[key] = true
	}
}

func TestMaxDegreeConstraint(t *testing.T) {

	n := 100

	_, edges, constraints := buildTestGraph(n, 0)

	deg := make([]int, n)

	for _, e := range edges {
		deg[e.U]++
		deg[e.V]++
	}

	for i := 0; i < n; i++ {

		if deg[i] > constraints.MaxDegree {
			t.Fatalf(
				"vertex %d degree %d exceeds maxDegree %d",
				i,
				deg[i],
				constraints.MaxDegree,
			)
		}
	}
}

func TestDistanceConstraint(t *testing.T) {

	points, edges, constraints := buildTestGraph(100, 0)

	for _, e := range edges {

		d := euclid(points[e.U], points[e.V])

		if d < constraints.DMin || d > constraints.DMax {
			t.Fatalf(
				"edge distance %f out of range [%f, %f]",
				d,
				constraints.DMin,
				constraints.DMax,
			)
		}
	}
}

func TestIsolatedVertices(t *testing.T) {

	n := 100
	isolated := 10

	_, edges, _ := buildTestGraph(n, isolated)

	deg := make([]int, n)

	for _, e := range edges {
		deg[e.U]++
		deg[e.V]++
	}

	isoCount := 0

	for i := 0; i < n; i++ {
		if deg[i] == 0 {
			isoCount++
		}
	}

	if isoCount < isolated {
		t.Fatalf(
			"expected at least %d isolated vertices, got %d",
			isolated,
			isoCount,
		)
	}
}
