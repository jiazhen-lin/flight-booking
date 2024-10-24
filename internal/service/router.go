package service

import (
	"container/heap"
	"time"

	"github.com/jiazhen-lin/flight-booking/internal/domain"
)

type path struct {
	airport         domain.AirportID
	flights         []domain.Flight
	arrivalTime     time.Time
	durationSeconds int64
}

// PriorityQueue implements a priority queue for flight paths
type PriorityQueue []path

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].durationSeconds != pq[j].durationSeconds {
		return pq[i].durationSeconds < pq[j].durationSeconds
	}
	return pq[i].arrivalTime.Before(pq[j].arrivalTime)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(path))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func findShortestNPaths(
	flights []domain.Flight, source, target domain.AirportID,
	departureTimeFrom, departureTimeTo time.Time, n int,
) ([]path, error) {
	// build the graph by adding available flights
	graph := make(map[domain.AirportID][]domain.Flight)
	for _, flight := range flights {
		graph[flight.DepartureAirportID] = append(graph[flight.DepartureAirportID], flight)
	}

	// using priority queue to implement dijkstra algorithm
	pq := &PriorityQueue{}
	heap.Push(pq, path{airport: source, arrivalTime: departureTimeFrom})
	visited := make(map[domain.AirportID]bool)

	var foundPaths []path
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(path)
		curAirport := cur.airport
		if visited[curAirport] {
			continue
		}

		// if the current path's arrival airport is the target, return the path
		if curAirport == target {
			// todo: update result
			foundPaths = append(foundPaths, cur)
			if len(foundPaths) >= n {
				return foundPaths, nil
			}
			continue
		}

		// check all connected flights
		for _, flight := range graph[curAirport] {
			// avoid loop
			if visited[flight.ArrivalAirportID] {
				continue
			}
			// ignore the flight that has already been departed
			if flight.DepartureTime.Before(cur.arrivalTime) {
				continue
			}
			// ignore the source flight that after the departure time limitation
			if curAirport == source && !flight.DepartureTime.Before(departureTimeTo) {
				continue
			}

			newArrivalTime := flight.DepartureTime.Add(time.Second * time.Duration(flight.DurationSeconds))
			newDuration := flight.DurationSeconds
			if len(cur.flights) > 0 {
				newDuration = int64(newArrivalTime.Sub(cur.flights[0].DepartureTime).Seconds())
			}
			cpy := make([]domain.Flight, len(cur.flights))
			copy(cpy, cur.flights)
			heap.Push(pq, path{
				airport:         flight.ArrivalAirportID,
				flights:         append(cpy, flight),
				arrivalTime:     newArrivalTime,
				durationSeconds: newDuration,
			})
		}

		visited[cur.airport] = true
	}

	return foundPaths, nil
}
