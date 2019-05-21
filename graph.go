package pathfinder

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

var (
	ErrNoPath = errors.New("no path to wanted destination")
)

type Node struct {
	Payload  interface{}
	Disabled bool
	Edges    []Edge
}

type Edge struct {
	Node *Node
	Cost float64
}

type graph struct {
	Nodes []*Node
}

func NewGraph() *graph {
	g := new(graph)
	return g
}

func (g *graph) nodeByPayload(payload interface{}) *Node {
	for _, n := range g.Nodes {
		if payload == n.Payload {
			return n
		}
	}
	panic(fmt.Sprintf("no vertex found with requested payload: %v", payload))
}

func (g *graph) AddNode(payload interface{}) {
	g.Nodes = append(g.Nodes, &Node{
		Payload: payload,
	})
}

func (g *graph) AddEdge(fromPayload, toPayload interface{}, cost float64) {
	fromNode := g.nodeByPayload(fromPayload)
	fromNode.Edges = append(fromNode.Edges, Edge{
		Node: g.nodeByPayload(toPayload),
		Cost: cost,
	})
}

func (g *graph) AddBidirectionalEdge(fromPayload, toPayload interface{}, cost float64) {
	g.AddEdge(fromPayload, toPayload, cost)
	g.AddEdge(toPayload, fromPayload, cost)
}

func (g *graph) EnableNode(payload interface{}) {
	g.nodeByPayload(payload).Disabled = false
}

func (g *graph) DisableNode(payload interface{}) {
	g.nodeByPayload(payload).Disabled = true
}

func (g *graph) ShortestPath(fromPayload, toPayload interface{}) ([]interface{}, float64, error) {
	var (
		from = g.nodeByPayload(fromPayload)
		to   = g.nodeByPayload(toPayload)
	)
	if from.Disabled {
		return []interface{}{}, math.Inf(1), ErrNoPath
	}

	distance := map[*Node]float64{}
	previousNode := map[*Node]*Node{}

	nodesLeft := make([]*Node, 0)

	for _, n := range g.Nodes {
		distance[n] = math.Inf(1)
		nodesLeft = append(nodesLeft, n)
	}
	distance[from] = 0

	for len(nodesLeft) > 0 {
		sort.Slice(nodesLeft, func(i, j int) bool {
			return distance[nodesLeft[i]] < distance[nodesLeft[j]]
		})
		chosenOne := nodesLeft[0]
		nodesLeft = append(nodesLeft[:0], nodesLeft[0+1:]...)

		if chosenOne == to {
			break
		}

		for _, edge := range chosenOne.Edges {
			if edge.Node.Disabled {
				continue
			}
			cost := distance[chosenOne] + edge.Cost
			if cost < distance[edge.Node] {
				distance[edge.Node] = cost
				previousNode[edge.Node] = chosenOne
			}
		}
	}

	// Traverse back to source to compute path
	path := make([]interface{}, 0)
	// Node we're currently looking at
	currentNode := to
	if previousNode[to] != nil || from == to {
		for currentNode != nil {
			path = append([]interface{}{currentNode.Payload}, path...)
			currentNode = previousNode[currentNode]
		}
	}

	if len(path) == 0 {
		return path, distance[to], ErrNoPath
	}

	return path, distance[to], nil
}
