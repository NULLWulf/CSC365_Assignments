package main

type Graph struct {
	Nodes []*BizNode
}

type BizNode struct {
	Name string
	Next []*BizNode
}