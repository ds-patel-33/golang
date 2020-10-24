package main

import (
	"container/list"
	"fmt"
)

type node struct {
	Id      string
	friends map[string]*node
}

func Nodes(n *node) []*node {

	visited := make(map[string]*node)

	queue := list.New()
	queue.PushBack(n)

	visited[n.Id] = n

	for queue.Len() > 0 {
		qnode := queue.Front()

		for id, node := range qnode.Value.(*node).friends {
			if _, ok := visited[id]; !ok {
				visited[id] = node
				queue.PushBack(node)
			}
		}
		queue.Remove(qnode)
	}

	nodes := make([]*node, 0)

	for _, node := range visited {
		nodes = append(nodes, node)
	}

	return nodes
}

func main() {
	node1 := &node{Id: "node1", friends: make(map[string]*node)}
	node2 := &node{Id: "node2", friends: make(map[string]*node)}
	node3 := &node{Id: "node3", friends: make(map[string]*node)}

	node2.friends[node1.Id] = node1

	node3.friends[node2.Id] = node2

	node1.friends[node3.Id] = node3

	root := &node{Id: "root", friends: make(map[string]*node)}
	n := make(map[string]*node)
	n[root.Id] = root
	n[root.Id].friends[node1.Id] = node1
	n[root.Id].friends[node2.Id] = node2
	n[root.Id].friends[node3.Id] = node3

	nodes := Nodes(root)
	for _, node := range nodes {
		fmt.Printf("%+v\n", node.Id)
	}
}
