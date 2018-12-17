package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func readInput(f string) []byte {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

type node struct {
	step     byte
	parents  []*node
	children []*node
}

type nodeSlice []*node

func (p nodeSlice) Len() int {
	return len(p)
}

func (p nodeSlice) Less(i, j int) bool {
	a, b := p[i].step, p[j].step
	return string(a) < string(b)
}

func (p nodeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func parseInstruction(s string) (byte, byte) {
	// s is the input string, ex., "Step D must be finished before step L can begin."
	return s[5], s[36]
}

func instructionGraph(instructions []string) map[byte]*node {
	nodes := make(map[byte]*node)

	for _, instruction := range instructions {
		a, b := parseInstruction(instruction)

		var ok bool
		var ante *node
		var post *node

		// Ensure that the ante and post nodes exist.
		if ante, ok = nodes[a]; !ok {
			ante = &node{
				step:     a,
				parents:  make([]*node, 0, 4),
				children: make([]*node, 0, 4),
			}
			nodes[a] = ante
		}
		if post, ok = nodes[b]; !ok {
			post = &node{
				step:     b,
				parents:  make([]*node, 0, 4),
				children: make([]*node, 0, 4),
			}
			nodes[b] = post
		}

		// Make the post node a child of the ante node.
		ante.children = append(ante.children, post)
		post.parents = append(post.parents, ante)
	}

	return nodes
}

func main() {
	input := readInput("../../input.txt")
	instructions := strings.Split(string(input), "\n")

	var graph = instructionGraph(instructions)

	var curr *node
	var queue nodeSlice
	var queued = make(map[*node]bool)
	var complete = make(map[*node]bool)

	for len(complete) < len(graph) {
		if len(queued) < len(graph) {
			for _, n := range graph {
				if queued[n] {
					continue
				}

				ready := true
				for _, p := range n.parents {
					ready = ready && complete[p]
				}

				if ready || len(n.parents) == 0 {
					queue = append(queue, n)
					queued[n] = true
				}
			}
		}

		sort.Sort(queue)
		curr, queue = queue[0], queue[1:]
		complete[curr] = true

		fmt.Printf("%s", string(curr.step))
	}
	fmt.Println()
}
