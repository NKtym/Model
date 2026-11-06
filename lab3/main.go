package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os/exec"
	"time"
)

type Point struct {
	X, Y float64
}

type Edge struct {
	U, V int
	D    float64
}

type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) bool {
	pa := d.Find(a)
	pb := d.Find(b)
	if pa == pb {
		return false
	}
	if d.r[pa] < d.r[pb] {
		d.p[pa] = pb
	} else if d.r[pb] < d.r[pa] {
		d.p[pb] = pa
	} else {
		d.p[pb] = pa
		d.r[pa]++
	}
	return true
}

type Constraints struct {
	UseMaxDeg          bool
	MaxDegree          int
	UseDistRange       bool
	DMin, DMax         float64
	UseMinDegree       bool
	MinDegree          int
	EnsureAllConnected bool
	NumIsolated        int
}

func randFloat(a, b float64) float64 {
	return a + rand.Float64()*(b-a)
}

func euclid(a, b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func normDist(d, W, H float64) float64 {
	return d / math.Hypot(W, H)
}

func probFromModel(dNorm, a, b float64, model int) float64 {
	if model == 1 {
		return math.Exp(-a * math.Pow(dNorm, b))
	}
	return 1.0 / (1.0 + math.Pow(dNorm, b))
}

func generatePoints(n int, W, H float64) []Point {
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i] = Point{X: rand.Float64() * W, Y: rand.Float64() * H}
	}
	return points
}

func buildGraph(points []Point, a, b float64, model int, constraints Constraints, seed int64) ([]Edge, int, int) {
	n := len(points)
	W, H := 100.0, 100.0
	r := rand.New(rand.NewSource(seed))
	dsu := NewDSU(n)
	edges := make([]Edge, 0, n)
	degree := make([]int, n)
	exists := make(map[[2]int]bool)
	isolated := make(map[int]bool)
	for len(isolated) < constraints.NumIsolated {
		v := r.Intn(n)
		isolated[v] = true
	}
	minEdges := 60
	maxEdges := n
	targetEdges := r.Intn(maxEdges-minEdges+1) + minEdges
	var startVertex, lastVertex int
	first := true

	for len(edges) < targetEdges {
		type candT struct {
			u, v int
			d    float64
			p    float64
		}
		cands := []candT{}
		for u := 0; u < n; u++ {
			if isolated[u] {
				continue
			}
			for v := u + 1; v < n; v++ {
				if isolated[v] {
					continue
				}
				key := [2]int{u, v}
				if exists[key] {
					continue
				}
				if constraints.UseMaxDeg && (degree[u] >= constraints.MaxDegree || degree[v] >= constraints.MaxDegree) {
					continue
				}
				d := euclid(points[u], points[v])
				if constraints.UseDistRange && (d < constraints.DMin || d > constraints.DMax) {
					continue
				}
				dNorm := normDist(d, W, H)
				p := probFromModel(dNorm, a, b, model)
				cands = append(cands, candT{u: u, v: v, d: d, p: p})
			}
		}
		if len(cands) == 0 {
			break
		}
		var chosen int
		sum := 0.0
		for _, c := range cands {
			sum += c.p
		}
		if sum == 0 {
			chosen = r.Intn(len(cands))
		} else {
			rv := r.Float64() * sum
			acc := 0.0
			for i, c := range cands {
				acc += c.p
				if rv <= acc {
					chosen = i
					break
				}
			}
		}
		c := cands[chosen]
		if degree[c.u] >= constraints.MaxDegree || degree[c.v] >= constraints.MaxDegree {
			continue
		}
		if first {
			startVertex = c.u
			first = false
		}
		lastVertex = c.v
		dsu.Union(c.u, c.v)
		key := [2]int{c.u, c.v}
		if c.u > c.v {
			key = [2]int{c.v, c.u}
		}
		exists[key] = true
		w := c.d
		if w > 100 {
			w = 100
		}
		edges = append(edges, Edge{U: c.u, V: c.v, D: w})
		degree[c.u]++
		degree[c.v]++
	}

	for u := 0; u < n; u++ {
		if isolated[u] || degree[u] > 0 {
			continue
		}
		for v := 0; v < n; v++ {
			if isolated[v] || u == v {
				continue
			}
			if degree[v] >= constraints.MaxDegree {
				continue
			}
			d := euclid(points[u], points[v])
			if constraints.UseDistRange && (d < constraints.DMin || d > constraints.DMax) {
				continue
			}
			key := [2]int{u, v}
			if u > v {
				key = [2]int{v, u}
			}
			if exists[key] {
				continue
			}
			dsu.Union(u, v)
			exists[key] = true
			w := d
			if w > 100 {
				w = 100
			}
			edges = append(edges, Edge{U: u, V: v, D: w})
			degree[u]++
			degree[v]++
			lastVertex = v
			break
		}
	}

	return edges, startVertex, lastVertex
}

func buildAdj(n int, edges []Edge) [][]int {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.U] = append(adj[e.U], e.V)
		adj[e.V] = append(adj[e.V], e.U)
	}
	return adj
}

func bfs(start int, adj [][]int) []int {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	queue := []int{start}
	dist[start] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue = append(queue, to)
			}
		}
	}
	return dist
}

func diameter(n int, edges []Edge) int {
	adj := buildAdj(n, edges)
	maxD := 0
	for i := 0; i < n; i++ {
		d := bfs(i, adj)
		for _, v := range d {
			if v > maxD {
				maxD = v
			}
		}
	}
	return maxD
}

func drawGraph(n int, points []Point, edges []Edge, filename string) {
	dot := "graph G {\n  layout=neato;\n  overlap=false;\n  splines=true;\n  node [shape=circle, fontsize=9, width=0.16];\n"
	for i := 0; i < n; i++ {
		dot += fmt.Sprintf("  \"V%d\" [pos=\"%f,%f!\"];\n", i, points[i].X, points[i].Y)
	}
	for _, e := range edges {
		dot += fmt.Sprintf("  \"V%d\" -- \"V%d\";\n", e.U, e.V)
	}
	dot += "}\n"
	tmp := "graph_tmp.dot"
	_ = ioutil.WriteFile(tmp, []byte(dot), 0644)
	_ = exec.Command("neato", "-Tpng", tmp, "-o", filename).Run()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	n := 100
	W, H := 100.0, 100.0
	constraints := Constraints{
		UseMaxDeg:    true,
		MaxDegree:    5,
		UseDistRange: true,
		DMin:         5,
		DMax:         100,
		UseMinDegree: true,
		MinDegree:    5,
		NumIsolated:  40,
	}
	points := generatePoints(n, W, H)
	models := []int{1, 2}
	for _, model := range models {
		seed := time.Now().UnixNano() + int64(model*1000)
		a := randFloat(0.02, 0.15)
		b := randFloat(1.1, 2.6)
		fmt.Printf("=== model %d ===\n", model)
		fmt.Printf("a=%.4f, b=%.4f\n", a, b)
		edges, startV, lastV := buildGraph(points, a, b, model, constraints, seed)
		fmt.Printf("Генерация началась с V%d, завершилась на V%d\n", startV, lastV)
		fmt.Printf("Edges: %d\n", len(edges))
		diam := diameter(n, edges)
		fmt.Printf("Diameter: %d\n", diam)
		outFile := fmt.Sprintf("graph_model_%d.png", model)
		drawGraph(n, points, edges, outFile)
		fmt.Printf("Saved: %s\n\n", outFile)
	}
}
