package main

import (
	"fmt"
)

type Vertex struct {
	visited    bool
	value      string
	neighbours []*Vertex
}

func NewVertex(value string) *Vertex {
	return &Vertex{
		value: value,

		visited:    false,
		neighbours: nil,
	}
}

func (v *Vertex) connect(vertex ...*Vertex) {
	v.neighbours = append(v.neighbours, vertex...)
}

type Graph struct{}

func (g *Graph) dfs(vertex *Vertex) {
	if vertex.visited {
		return
	}
	vertex.visited = true
	fmt.Println(vertex)
	for _, v := range vertex.neighbours {
		g.dfs(v)
	}
}

func (g *Graph) disconnected(vertices ...*Vertex) {
	for _, v := range vertices {
		g.dfs(v)
	}
}

func main() {
	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")
	v4 := NewVertex("D")
	v5 := NewVertex("E")
	g := Graph{}
	v1.connect(v2)
	v2.connect(v4, v5)
	v3.connect(v4, v5)
	g.dfs(v1)
}
