package pqueue

// import (
// 	"github.com/logpost/jobs-optimization-service/models"
// )

// An Item is something we manage in a priority queue.
type Item struct {
	// Job 				*models.Job
	Profit				float64
	JobID				string
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