package main

import (
	"container/heap"
	"encoding/gob"
	"log"
	"os"
)

type Graph struct {
	Nodes map[int]*gNode
}

type gNode struct {
	Key   int
	Value interface{}
	Edges []*Edge
	Root  int
}

type Edge struct {
	To     *gNode
	Weight float64
}

func (g *Graph) AddNode(key int, rkey int, value interface{}) {
	node := &gNode{
		Key:   key,
		Value: value,
		Root:  rkey,
		Edges: []*Edge{},
	}
	g.Nodes[key] = node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[int]*gNode),
	}
}

func (g *Graph) AddEdge(from, to int, weight float64) {
	fromNode, fromOK := g.Nodes[from]
	toNode, toOK := g.Nodes[to]
	if !fromOK || !toOK {
		return
	}
	edge := &Edge{
		To:     toNode,
		Weight: weight,
	}
	fromNode.Edges = append(fromNode.Edges, edge)
}

func (g *Graph) PrintGraph() {
	log.Println("Nodes:")
	for _, node := range g.Nodes {
		log.Printf("Key: %d, Value: %v\n", node.Key, node.Value)
		log.Println("Edges:")
		for _, edge := range node.Edges {
			log.Printf("  To: %d, Weight: %f\n", edge.To.Key, edge.Weight)
		}
	}
}

type nodeDist struct {
	nodeKey int
	dist    float64
}

type nodeHeap []*nodeDist

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h nodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*nodeDist))
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func (g *Graph) DijkstraShortestPath(startKey, targetKey int) ([]int, float64) {
	// Create a min-heap to store nodes with their minimum distance
	minHeap := &nodeHeap{}
	heap.Init(minHeap)

	// Create a map to store the minimum distance to each node
	distances := make(map[int]float64)

	// Create a map to store the previous node in the shortest path
	prevNodes := make(map[int]int)

	// Initialize distances to "infinity" for all nodes except the start node
	for key := range g.Nodes {
		if key == startKey {
			distances[key] = 0
		} else {
			distances[key] = 1e18 // a large number representing "infinity"
		}
	}

	// Add the start node to the min-heap
	heap.Push(minHeap, &nodeDist{nodeKey: startKey, dist: 0})

	// Dijkstra's algorithm
	for minHeap.Len() > 0 {
		current := heap.Pop(minHeap).(*nodeDist)
		currentNode := g.Nodes[current.nodeKey]

		// Process each adjacent node
		for _, edge := range currentNode.Edges {
			newDist := distances[current.nodeKey] + edge.Weight

			// If a shorter path is found, update the distance and previous node
			if newDist < distances[edge.To.Key] {
				distances[edge.To.Key] = newDist
				prevNodes[edge.To.Key] = current.nodeKey
				heap.Push(minHeap, &nodeDist{nodeKey: edge.To.Key, dist: newDist})
			}
		}
	}

	// Reconstruct the shortest path from start node to target node
	path := []int{}
	nodeKey := targetKey

	for nodeKey != startKey {
		path = append([]int{nodeKey}, path...)
		nodeKey = prevNodes[nodeKey]
	}

	path = append([]int{startKey}, path...)

	return path, distances[targetKey]
}

// SaveGraph the grab as a gob
func (g *Graph) SaveGraph() {
	// save the graph as a gob
	gob.Register(BusinessDataPoint{})
	f, err := os.Create("graph.gob")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(g)
	if err != nil {
		log.Fatal("encode error: ", err)
	}
}

// LoadGraph the graph as a gob
func LoadGraph() *Graph {
	gob.Register(BusinessDataPoint{})
	f, err := os.Open("graph.gob")
	if err != nil {
		log.Fatal("open file: ", err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var g Graph
	err = dec.Decode(&g)
	if err != nil {
		log.Fatal("decode error: ", err)
	}

	return &g
}

func (g *Graph) find(parent []int, i int) int {
	if parent[i] == -1 {
		return i
	}
	return g.find(parent, parent[i])
}

func (g *Graph) Union() int {
	n := len(g.Nodes)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = -1
	}
	for _, node := range g.Nodes {
		root1 := g.find(parent, node.Key)
		for _, edge := range node.Edges {
			root2 := g.find(parent, edge.To.Key)
			if root1 != root2 {
				parent[root1] = root2
			}
		}
	}
	count := 0
	for i := 0; i < n; i++ {
		if parent[i] == -1 {
			count++
		}
	}
	return count
}
