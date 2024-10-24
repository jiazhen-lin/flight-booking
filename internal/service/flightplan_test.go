package service

import (
	"fmt"
	"testing"

	"github.com/RyanCarrier/dijkstra/v2"
	"github.com/dominikbraun/graph"
	"github.com/stretchr/testify/require"
)

func TestFlightPlanService_Search(t *testing.T) {
	type City struct {
		Name string
	}

	cityHash := func(c City) string {
		return c.Name
	}

	london := City{Name: "london"}
	munich := City{Name: "munich"}
	paris := City{Name: "paris"}
	madrid := City{Name: "madrid"}

	g := graph.New(cityHash, graph.Weighted())

	_ = g.AddVertex(london)
	_ = g.AddVertex(munich)
	_ = g.AddVertex(paris)
	_ = g.AddVertex(madrid)

	_ = g.AddEdge("london", "munich", graph.EdgeWeight(3))
	_ = g.AddEdge("london", "paris", graph.EdgeWeight(2))
	_ = g.AddEdge("london", "madrid", graph.EdgeWeight(5))
	_ = g.AddEdge("munich", "madrid", graph.EdgeWeight(6))
	_ = g.AddEdge("munich", "paris", graph.EdgeWeight(2))
	_ = g.AddEdge("paris", "madrid", graph.EdgeWeight(4))

	path, _ := graph.ShortestPath(g, "london", "madrid")
	fmt.Println(path)
}

func TestFlightPlanService_Search2(t *testing.T) {
	graph := dijkstra.NewGraph()
	//Add the 3 verticies
	graph.AddVertexAndArcs(0, map[int]uint64{1: 1, 2: 1})
	graph.AddVertexAndArcs(1, map[int]uint64{3: 1})
	graph.AddVertexAndArcs(2, map[int]uint64{3: 1})
	graph.AddVertexAndArcs(3, map[int]uint64{4: 1})

	best, err := graph.ShortestAll(0, 4)
	require.NoError(t, err)
	fmt.Println("Shortest distances are", best.Distance, "with paths; ", best.Paths)

	best, err = graph.LongestAll(0, 4)
	require.NoError(t, err)
	fmt.Println("Longest distances are", best.Distance, "following path ", best.Paths)
}
