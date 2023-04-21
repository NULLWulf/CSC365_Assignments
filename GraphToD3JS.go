package main

//
//import (
//	"encoding/json"
//	"strconv"
//)
//
//type D3Graph struct {
//	Nodes []Node `json:"nodes"`
//	Edges []Edge `json:"edges"`
//}
//
//type D3Node struct {
//	ID           string  `json:"id"`
//	Label        string  `json:"label"`
//	Group        string  `json:"group"`
//	BusinessData string  `json:"businessData"`
//	Value        float64 `json:"value"`
//}
//
//type D3Edge struct {
//	Source string  `json:"source"`
//	Target string  `json:"target"`
//	Weight float64 `json:"weight"`
//	Color  string  `json:"color"`
//}
//
//func ShortPathListToD3Graph(bizFileKeys []int) ([]byte, error) {
//	graph := D3Graph{
//		Nodes: []D3Node{},
//		Edges: []D3Edge{},
//	}
//	for _, key := range bizFileKeys {
//		node := D3Node{
//			ID:           strconv.Itoa(key),
//			Label:        strconv.Itoa(key),
//			Group:        "business",
//			BusinessData: strconv.Itoa(key),
//			Value:        1,
//		}
//		graph.Nodes = append(graph.Nodes, node)
//	}
//	for _, node := range graph.Nodes {
//		for _, edge := range node.Edges {
//			graph.Edges = append(graph.Edges, edge)
//		}
//	}
//	return json.Marshal(graph)
//}
