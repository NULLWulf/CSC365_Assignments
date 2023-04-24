package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
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

type tNode struct {
	Key   int
	Value interface{}
	Edges []*tEdge
	Root  int
}

type tEdge struct {
	ToKey  int
	Weight float64
}

type nodeDist struct {
	nodeKey int
	dist    float64
}

type nodeHeap []*nodeDist

// AddNode adds a node to the graph
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

// AddEdge adds an edge to the graph
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

func (g *Graph) DijkstraShortestPath(startKey, targetKey int) ([]int, float64, error) {
	// Create a min-heap to store nodes with their minimum distance
	minHeap := &nodeHeap{}
	heap.Init(minHeap)

	// Create a map to store the minimum distance to each node
	distances := make(map[int]float64)

	// Create a map to store the previous node in the shortest path
	prevNodes := make(map[int]int)

	// Initialize distances to "infinity" for all nodes except the start node
	for key := range g.Nodes {
		if key == startKey { // start node
			distances[key] = 0 // distance to itself is 0
		} else {
			distances[key] = 1e18 // a large number representing "infinity"
		}
	}

	// Add the start node to the min-heap
	heap.Push(minHeap, &nodeDist{nodeKey: startKey, dist: 0})

	// Dijkstra's algorithm
	for minHeap.Len() > 0 { // while there are nodes in the min-heap
		current := heap.Pop(minHeap).(*nodeDist) // pop the node with the shortest distance from the min-heap
		currentNode := g.Nodes[current.nodeKey]  // get the corresponding node

		for _, edge := range currentNode.Edges { // Process each adjacent node

			newDist := distances[current.nodeKey] + edge.Weight // Calculate the new distance

			if newDist < distances[edge.To.Key] { // If a shorter path is found, update the distance and previous node

				distances[edge.To.Key] = newDist                                   // update the distance
				prevNodes[edge.To.Key] = current.nodeKey                           // update the previous node
				heap.Push(minHeap, &nodeDist{nodeKey: edge.To.Key, dist: newDist}) // push the adjacent node to the min-heap
			}
		}
	}
	// Check if we've found a path to the target node
	if _, ok := prevNodes[targetKey]; !ok {
		return nil, 0, fmt.Errorf("target node %d not found in graph", targetKey)
	}

	// Reconstruct the shortest path from start node to target node
	path := []int{}      // initialize an empty path
	nodeKey := targetKey // start from the target node

	for nodeKey != startKey { // while we haven't reached the start node yet
		path = append([]int{nodeKey}, path...) // prepend the current node to the path
		nodeKey = prevNodes[nodeKey]           // move to the previous node
	}

	path = append([]int{startKey}, path...) // prepend the start node to the path

	return path, distances[targetKey], nil // return the path and its distance
}

func (g *Graph) find(parent []int, i int) int {
	if parent[i] == -1 {
		return i
	}
	return g.find(parent, parent[i])
}

// UnionFind returns the number of connected components in the graph
func (g *Graph) UnionFind() int {
	n := len(g.Nodes)
	parent := make([]int, n)
	for i := 0; i < n; i++ { // initialize parent array
		parent[i] = -1
	}
	for _, node := range g.Nodes { // iterate through all nodes
		root1 := g.find(parent, node.Key) // find the root of the current node
		for _, edge := range node.Edges { // iterate through all edges of the current node
			root2 := g.find(parent, edge.To.Key) // find the root of the adjacent node
			if root1 != root2 {                  // if the roots are different, union them
				parent[root1] = root2 // set the parent of root1 as root2
			}
		}
	}
	count := 0               // initialize the number of roots
	for i := 0; i < n; i++ { // count the number of roots
		if parent[i] == -1 { // if the parent is -1, it is a root
			count++
		}
	}
	return count
}

// serializes the graph into a JSON file
func (g *Graph) serialize() {
	nodes := make([]tNode, len(g.Nodes))
	for i, n := range g.Nodes {
		// temp node
		tn := tNode{
			Key:   n.Key,
			Value: n.Value,
			Edges: make([]*tEdge, len(n.Edges)),
			Root:  n.Root,
		}

		// iterate over each edge and create a temp edge for it
		for j, e := range n.Edges {
			te := &tEdge{
				ToKey:  e.To.Key,
				Weight: e.Weight,
			}
			tn.Edges[j] = te
		}

		// add the temp node to the slice of nodes
		nodes[i] = tn
	}

	// Serialize the slice of nodes
	b, err := json.Marshal(nodes)
	if err != nil {
		log.Fatal(err)
	}

	// Write the serialized nodes to a file
	err = os.WriteFile("graph.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node: %v\n", len(nodes))
}

// deserializeGraph reads the serialized graph from the file and returns a Graph object
func deserializeGraph() (*Graph, error) {
	// Read the contents of the serialized file
	b, err := os.ReadFile("graph.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal the serialized nodes into a slice of temp nodes
	var tnodes []tNode
	err = json.Unmarshal(b, &tnodes)
	if err != nil {
		return nil, err
	}

	// Create a new Graph object
	graph := &Graph{
		Nodes: make(map[int]*gNode),
	}

	// Create a gNode object for each temp node and add it to the Nodes map
	for _, tn := range tnodes {
		gn := &gNode{
			Key:   tn.Key,
			Value: tn.Value,
			Edges: make([]*Edge, len(tn.Edges)),
			Root:  tn.Root,
		}

		graph.Nodes[tn.Key] = gn
	}

	// Create Edge objects for each temp edge and add them to the Edges slice of the corresponding gNode
	for _, tn := range tnodes {
		from := graph.Nodes[tn.Key]

		for j, te := range tn.Edges {
			to := graph.Nodes[te.ToKey]
			e := &Edge{
				To:     to,
				Weight: te.Weight,
			}
			from.Edges[j] = e
		}
	}

	return graph, nil
}

// getRandomNodes returns n random nodes from the graph
func (g *Graph) getRandomNodes(n int) []int {
	rand.Seed(time.Now().UnixNano())
	points := make([]int, n)
	dups := make(map[int]bool)
	for i := 0; i < n; i++ {
		if _, ok := dups[i]; ok {
			continue
		} else {
			points[i] = rand.Intn(len(g.Nodes))
		}
	}
	return points
}
