package pqueue

import (
	// "container/heap"
	"github.com/logpost/jobs-optimization-service/models"
)

// An Item is something we manage in a priority queue.
type Item struct {
	Job 				*models.Job
	Profit				float64
	JobID				string
	DistanceToOrigin	float64
	index				int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {

	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Profit > pq[j].Profit
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j]	=	pq[j], pq[i]
	pq[i].index		=	i
	pq[j].index		=	j
}

// Push do push item in to queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n			:=	len(*pq)
	item		:=	x.(*Item)
	item.index	=	n
	*pq			=	append(*pq, item)
}

// Pop do pop item in to queue.
func (pq *PriorityQueue) Pop() interface{} {
	old			:=	*pq
	n			:=	len(old)
	item		:=	old[n-1]
	old[n-1]	=	nil 	// avoid memory leak
	item.index	=	-1		// for safety
	*pq			=	old[0 : n-1]

	return item
}

// update modifies the priority and value of an Item in the queue.
// func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	
// 	item.value		=	value
// 	item.priority	=	priority

// 	heap.Fix(pq, item.index)
// }

// // This example creates a PriorityQueue with some items, adds and manipulates an item,
// // and then removes the items in priority order.
// func main() {
// 	// Some items and their priorities.
// 	items := map[string]int{
// 		"banana": 3, "apple": 2, "pear": 4,
// 	}

// 	// Create a priority queue, put the items in it, and
// 	// establish the priority queue (heap) invariants.
// 	pq := make(PriorityQueue, len(items))
// 	i := 0
// 	for value, priority := range items {
// 		pq[i] = &Item{
// 			value:    value,
// 			priority: priority,
// 			index:    i,
// 		}
// 		i++
// 	}
// 	heap.Init(&pq)

// 	// Insert a new item and then modify its priority.
// 	item := &Item{
// 		value:    "orange",
// 		priority: 1,
// 	}
// 	heap.Push(&pq, item)
// 	pq.update(item, item.value, 5)

// 	// Take the items out; they arrive in decreasing priority order.
// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(*Item)
// 		fmt.Printf("%.2d:%s ", item.priority, item.value)
// 	}
// }