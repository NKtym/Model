package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os/exec"
	"sort"
)

type Edge struct {
	U, V string
	W    int
}

type EdgeInfo struct {
	U, V     string
	tij      int
	t_rn     int
	t_ro     int
	t_pn     int
	t_po     int
	Rij      int
	rij      int
	Critical bool
}

func main() {
	edges := []Edge{
		{"A", "B", 4},
		{"A", "C", 6},
		{"B", "D", 3},
		{"C", "E", 5},
		{"D", "F", 4},
		{"E", "G", 3},
		{"F", "G", 3},
	}
	nodesMap := map[string]bool{}
	for _, e := range edges {
		nodesMap[e.U] = true
		nodesMap[e.V] = true
	}
	nodesMap["A"] = true
	nodes := make([]string, 0, len(nodesMap))
	for k := range nodesMap {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)
	idx := map[string]int{}
	for i, n := range nodes {
		idx[n] = i
	}
	n := len(nodes)
	adj := make([][]Edge, n)
	indeg := make([]int, n)
	outdeg := make([]int, n)
	for _, e := range edges {
		u := idx[e.U]
		v := idx[e.V]
		adj[u] = append(adj[u], e)
		indeg[v]++
		outdeg[u]++
	}
	q := []int{}
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	toporder := []int{}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		toporder = append(toporder, u)
		for _, e := range adj[u] {
			v := idx[e.V]
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	if len(toporder) != n {
		fmt.Println("Ошибка: граф содержит цикл")
		return
	}
	VE := make([]int, n)
	for i := 0; i < n; i++ {
		VE[i] = math.MinInt32
	}
	hasIncoming := map[int]bool{}
	for _, e := range edges {
		hasIncoming[idx[e.V]] = true
	}
	for i := 0; i < n; i++ {
		if !hasIncoming[i] {
			VE[i] = 0
		}
	}
	for _, u := range toporder {
		if VE[u] == math.MinInt32 {
			continue
		}
		for _, e := range adj[u] {
			v := idx[e.V]
			if VE[v] < VE[u]+e.W {
				VE[v] = VE[u] + e.W
			}
		}
	}
	T := math.MinInt32
	for i := 0; i < n; i++ {
		if VE[i] != math.MinInt32 && VE[i] > T {
			T = VE[i]
		}
	}
	VL := make([]int, n)
	for i := range VL {
		VL[i] = math.MaxInt32
	}
	for i := 0; i < n; i++ {
		if outdeg[i] == 0 && VE[i] != math.MinInt32 {
			if VL[i] > VE[i] {
				VL[i] = VE[i]
			}
		}
	}
	for i := len(toporder) - 1; i >= 0; i-- {
		u := toporder[i]
		for _, e := range adj[u] {
			v := idx[e.V]
			if VL[v] == math.MaxInt32 {
				continue
			}
			candidate := VL[v] - e.W
			if candidate < VL[u] {
				VL[u] = candidate
			}
		}
	}
	for i := 0; i < n; i++ {
		if VL[i] == math.MaxInt32 && VE[i] != math.MinInt32 {
			VL[i] = VE[i]
		}
	}
	edgeInfos := []EdgeInfo{}
	for _, e := range edges {
		u := idx[e.U]
		v := idx[e.V]
		tij := e.W
		t_rn := VE[u]
		t_ro := t_rn + tij
		t_po := VL[v]
		t_pn := -1
		if t_po != math.MaxInt32 {
			t_pn = t_po - tij
		}
		if t_rn == math.MinInt32 {
			t_rn = -1
			t_ro = -1
		}
		if t_po == math.MaxInt32 {
			t_po = -1
			t_pn = -1
		}
		Rij := -1
		if t_pn >= 0 && t_rn >= 0 {
			Rij = t_pn - t_rn
		}
		ei := EdgeInfo{
			U:        e.U,
			V:        e.V,
			tij:      tij,
			t_rn:     t_rn,
			t_ro:     t_ro,
			t_pn:     t_pn,
			t_po:     t_po,
			Rij:      Rij,
			rij:      0,
			Critical: (Rij == 0),
		}
		edgeInfos = append(edgeInfos, ei)
	}
	for i := range edgeInfos {
		ei := &edgeInfos[i]
		vIdx := idx[ei.V]
		var minSuccStart int
		if outdeg[vIdx] == 0 {
			minSuccStart = T
		} else {
			if VE[vIdx] != math.MinInt32 {
				minSuccStart = VE[vIdx]
			} else {
				minSuccStart = T
			}
		}
		if ei.t_ro >= 0 {
			free := minSuccStart - ei.t_ro
			if free < 0 {
				free = 0
			}
			ei.rij = free
		} else {
			ei.rij = 0
		}
	}
	critAdj := map[string][]string{}
	starts := []string{}
	ends := map[string]bool{}
	for _, ei := range edgeInfos {
		if ei.Critical {
			critAdj[ei.U] = append(critAdj[ei.U], ei.V)
		}
	}
	for i := 0; i < n; i++ {
		if !hasIncoming[i] {
			starts = append(starts, nodes[i])
		}
		if outdeg[i] == 0 {
			ends[nodes[i]] = true
		}
	}
	critPaths := [][]string{}
	var stack []string
	var dfs func(string)
	dfs = func(u string) {
		stack = append(stack, u)
		if ends[u] {
			if len(stack) > 1 {
				cp := make([]string, len(stack))
				copy(cp, stack)
				critPaths = append(critPaths, cp)
			}
		}
		for _, v := range critAdj[u] {
			dfs(v)
		}
		stack = stack[:len(stack)-1]
	}
	for _, s := range starts {
		if _, ok := critAdj[s]; ok {
			dfs(s)
		}
	}
	fmt.Printf("%-6s %-5s %-7s %-7s %-7s %-7s %-5s %-5s\n", "Шифр", "tij", "t_rn", "t_ro", "t_pn", "t_po", "Rij", "rij")
	for _, ei := range edgeInfos {
		t_rn_s := "-"
		t_ro_s := "-"
		t_pn_s := "-"
		t_po_s := "-"
		if ei.t_rn >= 0 {
			t_rn_s = fmt.Sprintf("%d", ei.t_rn)
		}
		if ei.t_ro >= 0 {
			t_ro_s = fmt.Sprintf("%d", ei.t_ro)
		}
		if ei.t_pn >= 0 {
			t_pn_s = fmt.Sprintf("%d", ei.t_pn)
		}
		if ei.t_po >= 0 {
			t_po_s = fmt.Sprintf("%d", ei.t_po)
		}
		Rij_s := ei.Rij
		rij_s := ei.rij
		fmt.Printf("%-6s %-5d %-7s %-7s %-7s %-7s %-5d %-5d", ei.U+"-"+ei.V, ei.tij, t_rn_s, t_ro_s, t_pn_s, t_po_s, Rij_s, rij_s)
		fmt.Println()
	}
	if len(critPaths) > 0 {
		fmt.Println("\nКритические пути:")
		for _, p := range critPaths {
			fmt.Println(" ", p)
		}
	} else {
		fmt.Println("\nКритические пути не найдены.")
	}
	dot := "digraph G {\n"
	dot += "  rankdir=LR;\n"
	dot += "  node [shape=circle, style=filled, fillcolor=white, fontsize=12];\n"
	dot += "  edge [fontsize=10];\n\n"
	for _, name := range nodes {
		dot += fmt.Sprintf("  \"%s\" [label=\"%s\"];\n", name, name)
	}
	dot += "\n"
	for _, ei := range edgeInfos {
		label := fmt.Sprintf("%d", ei.tij)
		color := "gray"
		penwidth := 1.0
		if ei.Critical {
			color = "red"
			penwidth = 2.5
		}
		dot += fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\", color=%s, penwidth=%.1f];\n", ei.U, ei.V, label, color, penwidth)
	}
	dot += "}\n"
	if err := ioutil.WriteFile("graph.dot", []byte(dot), 0644); err != nil {
		fmt.Println("Ошибка записи graph.dot:", err)
		return
	}
	fmt.Println("\nФайл graph.dot создан.")
	cmd := exec.Command("dot", "-Tpng", "graph.dot", "-o", "graph.png")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ошибка при выполнении dot: %v\nВывод: %s\n", err, string(out))
		return
	}
	fmt.Println("graph.png успешно создан.")
	fmt.Println("T_кр =", T)
}
