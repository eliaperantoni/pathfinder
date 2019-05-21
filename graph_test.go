package pathfinder_test

import (
	"fmt"
	"github.com/eliaperantoni/pathfinder"
	"math"
	"math/rand"
	"testing"
)

func TestGraph_AddNode(t *testing.T) {
	g := pathfinder.NewGraph()
	nodes := []interface{}{"Go", "Is", "Cool"}
	for i, n := range nodes {
		g.AddNode(n)
		got := g.Nodes[i].Payload.(string)
		want := n.(string)
		if got != want {
			t.Errorf("g.Nodes[%d] = %s; want %s", i, got, want)
		}
	}
}

func TestGraph_AddEdge(t *testing.T) {
	g := pathfinder.NewGraph()
	nodes := []interface{}{"a", "b", "c", "d"}
	edges := map[interface{}][]interface{}{
		"a": {
			"b",
			"c",
			"d",
		},
		"b": {
			"c",
			"d",
		},
		"d": {
			"c",
		},
	}
	for _, n := range nodes {
		g.AddNode(n)
	}
	for from, edges := range edges {
		for _, to := range edges {
			g.AddEdge(from, to, 1)
		}
	}
}

func TestGraph_ShortestPath(t *testing.T) {
	type shortestPathOutput struct {
		path []interface{}
		cost float64
		err  error
	}

	checkShortestPathOutput := func(t *testing.T, got shortestPathOutput, want shortestPathOutput) {
		if got.err != want.err {
			t.Errorf("graph.Shortestpath() err = %v, want %v", got.err, want.err)
		}
		if got.cost != want.cost {
			t.Errorf("graph.Shortestpath() cost = %.2f, want %.2f", got.cost, want.cost)
		}
		if len(got.path) != len(want.path) {
			t.Errorf("graph.Shortestpath() path len = %d, want %d", len(got.path), len(want.path))
		}
		for i := 0; i < len(want.path); i++ {
			if i+1 > len(got.path) {
				break
			}
			if got.path[i] != want.path[i] {
				t.Errorf("graph.Shortestpath()[%d] = %v, want %v", i, got.path[i], want.path[i])
			}
		}
	}

	type _node struct {
		payload  interface{}
		disabled bool
	}
	type _edge struct {
		node          interface{}
		cost          float64
		bidirectional bool
	}
	tests := map[string]struct {
		nodes []_node
		edges map[interface{}][]_edge
		from  interface{}
		to    interface{}
		want  shortestPathOutput
	}{
		"Very simple": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
			},
			"a",
			"b",
			shortestPathOutput{
				[]interface{}{"a", "b"},
				1,
				nil,
			},
		},
		"Bidirectional": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						12,
						true,
					},
				},
			},
			"b",
			"a",
			shortestPathOutput{
				[]interface{}{"b", "a"},
				12,
				nil,
			},
		},
		"Slightly more complex": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
				{payload: "c"},
				{payload: "d"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
					{
						"d",
						12,
						false,
					},
					{
						"c",
						3,
						false,
					},
				},
				"b": {
					{
						"c",
						2,
						false,
					},
					{
						"d",
						4,
						false,
					},
				},
				"d": {
					{
						"c",
						6,
						false,
					},
				},
			},
			"a",
			"d",
			shortestPathOutput{
				[]interface{}{"a", "b", "d"},
				5,
				nil,
			},
		},
		"No path": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
				{payload: "c"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
			},
			"a",
			"c",
			shortestPathOutput{
				[]interface{}{},
				math.Inf(1),
				pathfinder.ErrNoPath,
			},
		},
		"Disabled node, can still succeed": {
			[]_node{
				{payload: "a"},
				{payload: "b", disabled: true},
				{payload: "c"},
				{payload: "d"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
					{
						"c",
						10,
						false,
					},
				},
				"b": {
					{
						"d",
						1,
						false,
					},
				},
				"c": {
					{
						"d",
						1,
						false,
					},
				},
			},
			"a",
			"d",
			shortestPathOutput{
				[]interface{}{"a", "c", "d"},
				11,
				nil,
			},
		},
		"Disabled node, no path": {
			[]_node{
				{payload: "a"},
				{payload: "b", disabled: true},
				{payload: "c"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
				"b": {
					{
						"c",
						1,
						false,
					},
				},
			},
			"a",
			"c",
			shortestPathOutput{
				[]interface{}{},
				math.Inf(1),
				pathfinder.ErrNoPath,
			},
		},
		"Disabled node is DST": {
			[]_node{
				{payload: "a"},
				{payload: "b", disabled: true},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
			},
			"a",
			"b",
			shortestPathOutput{
				[]interface{}{},
				math.Inf(1),
				pathfinder.ErrNoPath,
			},
		},
		"SRC and DST match": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
			},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{"a"},
				0,
				nil,
			},
		},
		"SRC and DST match, no edges": {
			[]_node{
				{payload: "a"},
			},
			map[interface{}][]_edge{},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{"a"},
				0,
				nil,
			},
		},
		"SRC and DST match and is disabled": {
			[]_node{
				{payload: "a", disabled: true},
				{payload: "b"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
			},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{},
				math.Inf(1),
				pathfinder.ErrNoPath,
			},
		},
		"Loop": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
				{payload: "c"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
				"b": {
					{
						"c",
						1,
						false,
					},
				},
				"c": {
					{
						"a",
						1,
						false,
					},
				},
			},
			"a",
			"c",
			shortestPathOutput{
				[]interface{}{"a", "b", "c"},
				2,
				nil,
			},
		},
		"Loop, SRC and DST match": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
				{payload: "c"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						false,
					},
				},
				"b": {
					{
						"c",
						1,
						false,
					},
				},
				"c": {
					{
						"a",
						1,
						false,
					},
				},
			},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{"a"},
				0,
				nil,
			},
		},
		"Loop, bidirectional": {
			[]_node{
				{payload: "a"},
				{payload: "b"},
				{payload: "c"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"b",
						1,
						true,
					},
				},
				"b": {
					{
						"c",
						1,
						true,
					},
				},
				"c": {
					{
						"a",
						1,
						true,
					},
				},
			},
			"a",
			"c",
			shortestPathOutput{
				[]interface{}{"a", "c"},
				1,
				nil,
			},
		},
		"Node has edge to itself, cost 1": {
			[]_node{
				{payload: "a"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"a",
						1,
						true,
					},
				},
			},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{"a"},
				0,
				nil,
			},
		},
		"Node has edge to itself, cost 0": {
			[]_node{
				{payload: "a"},
			},
			map[interface{}][]_edge{
				"a": {
					{
						"a",
						0,
						true,
					},
				},
			},
			"a",
			"a",
			shortestPathOutput{
				[]interface{}{"a"},
				0,
				nil,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := pathfinder.NewGraph()
			for _, n := range tt.nodes {
				g.AddNode(n.payload)
				if n.disabled {
					g.DisableNode(n.payload)
				}
			}
			for from, edges := range tt.edges {
				for _, edge := range edges {
					g.AddEdge(from, edge.node, edge.cost)
					if edge.bidirectional {
						g.AddEdge(edge.node, from, edge.cost)
					}
				}
			}
			gotPath, gotCost, gotErr := g.ShortestPath(tt.from, tt.to)
			checkShortestPathOutput(t, shortestPathOutput{
				gotPath,
				gotCost,
				gotErr,
			}, tt.want)
		})
	}

	t.Run("toggle node", func(t *testing.T) {
		g := pathfinder.NewGraph()
		g.AddNode("a")
		g.AddNode("b")
		g.AddNode("c")
		g.AddNode("d")
		g.AddEdge("a", "b", 1)
		g.AddEdge("a", "c", 10)
		g.AddEdge("b", "d", 1)
		g.AddEdge("c", "d", 1)

		g.DisableNode("b")
		gotPath, gotCost, gotErr := g.ShortestPath("a", "d")
		checkShortestPathOutput(t, shortestPathOutput{
			gotPath,
			gotCost,
			gotErr,
		}, shortestPathOutput{
			[]interface{}{"a", "c", "d"},
			11,
			nil,
		})

		g.EnableNode("b")
		gotPath, gotCost, gotErr = g.ShortestPath("a", "d")
		checkShortestPathOutput(t, shortestPathOutput{
			gotPath,
			gotCost,
			gotErr,
		}, shortestPathOutput{
			[]interface{}{"a", "b", "d"},
			2,
			nil,
		})
	})
}

func BenchmarkGraph_ConstructFullyConnected100Nodes(b *testing.B) {
	const nodesN = 100

	for n := 0; n < b.N; n++ {
		g := pathfinder.NewGraph()
		for i := 0; i < nodesN; i++ {
			g.AddNode(i)
		}
		for from := 0; from < nodesN; from++ {
			for to := 0; to < nodesN; to++ {
				if from == to {
					continue
				}
				g.AddEdge(from, to, 1)
			}
		}
	}
}

func benchmarkGraph_ShortestPathXNodes(b *testing.B, nodesN int) {
	const dropoutRate = .3

	g := pathfinder.NewGraph()
	for i := 0; i < nodesN; i++ {
		g.AddNode(i)
	}
	for from := 0; from < nodesN; from++ {
		for to := 0; to < nodesN; to++ {
			if from == to {
				continue
			}
			if rand.Float64() < dropoutRate {
				continue
			}
			g.AddEdge(from, to, 1)
		}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.ShortestPath(0, nodesN-1)
	}
}

func BenchmarkGraph_ShortestPath100Nodes(b *testing.B) {
	benchmarkGraph_ShortestPathXNodes(b, 100)
}

func BenchmarkGraph_ShortestPath1000Nodes(b *testing.B) {
	benchmarkGraph_ShortestPathXNodes(b, 1000)
}

func ExampleGraph_ShortestPath() {
	graph := pathfinder.NewGraph()

	// In this example we're building this graph
	//
	//		A-------------------+
	//		 \                   \
	//        B                   C
	//         \                   \
	//          D-------------------+
	//
	// And finding the shortest path from A to D, assuming
	// taking the C route is much more costly

	// graph.AddNode(content)
	graph.AddNode("a")
	graph.AddNode("b")
	graph.AddNode("c")
	graph.AddNode("d")

	// graph.AddEdge(source, destination, cost)
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 10)
	graph.AddEdge("b", "d", 1)
	graph.AddEdge("c", "d", 10)

	path, cost, err := graph.ShortestPath("a", "d")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Path:")
	for _, node := range path {
		fmt.Printf("\t%s\n", node)
	}
	fmt.Printf("Cost: %.2f\n", cost)
	// Output: Path:
	//	a
	//	b
	//	d
	// Cost: 2.00
}
