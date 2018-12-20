package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
	duration uint32
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

func (n *node) setDuration() {
	n.duration = uint32(n.step - 4)
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
			ante.setDuration()
			nodes[a] = ante
		}
		if post, ok = nodes[b]; !ok {
			post = &node{
				step:     b,
				parents:  make([]*node, 0, 4),
				children: make([]*node, 0, 4),
			}
			post.setDuration()
			nodes[b] = post
		}

		// Make the post node a child of the ante node.
		ante.children = append(ante.children, post)
		post.parents = append(post.parents, ante)
	}

	return nodes
}

func getNodesNext(waiting nodeSlice, complete map[*node]bool) (nodeSlice, nodeSlice) {
	var rdy = make(nodeSlice, 0)
	var rem = make(nodeSlice, 0, len(waiting))

	for _, n := range waiting {
		ready := true
		for _, p := range n.parents {
			ready = ready && complete[p]
		}

		if ready || len(n.parents) == 0 {
			rdy = append(rdy, n)
		} else {
			rem = append(rem, n)
		}
	}

	return rdy, rem
}

func nextNode(ready, waiting nodeSlice, complete map[*node]bool) (*node, nodeSlice, nodeSlice) {
	var next *node

	r, waiting := getNodesNext(waiting, complete)
	ready = append(ready, r...)

	if len(ready) > 0 {
		sort.Sort(ready)
		next, ready = ready[0], ready[1:]
	}

	return next, ready, waiting
}

func main() {
	input := readInput("../../input.txt")
	instructions := strings.Split(string(input), "\n")

	var graph = instructionGraph(instructions)
	var waiting nodeSlice
	for _, n := range graph {
		waiting = append(waiting, n)
	}

	var cnt uint32
	var wip = [5]*node{}
	var ready nodeSlice
	var complete = make(map[*node]bool)

	for len(complete) < len(graph) {
		fmt.Printf("> %4d", cnt)
		for _, n := range wip {
			if n == nil {
				fmt.Printf(" - #")
			} else {
				fmt.Printf(" - %s", string(n.step))
			}
		}
		fmt.Println()

		// Fast Forward logic
		// Look for the smallest remaining time so that this time can be
		// subtracted from all in-progress steps.
		var t uint32 = math.MaxUint32
		for _, n := range wip {
			if n != nil {
				if n.duration < t {
					t = n.duration
				}
			}
		}
		if t != math.MaxUint32 {
			cnt += t
		} else {
			t = 0
		}

		//var inc = false
		var nxt *node
		for i, n := range wip {
			if n == nil {
				fmt.Printf(" --- #,  #\n")
				n, ready, waiting = nextNode(ready, waiting, complete)
				wip[i] = n
				continue
			}

			fmt.Printf(" --- %s, %2d", string(n.step), n.duration)
			n.duration -= t

			if n.duration == 0 {
				fmt.Printf(" --> %s", string(n.step))

				complete[n] = true
				nxt, ready, waiting = nextNode(ready, waiting, complete)
				wip[i] = nxt
			}
			fmt.Println()
		}
	}
	fmt.Printf("\nTotal time: %d\n", cnt)
}
