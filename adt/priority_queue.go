package adt

type priorityNode struct {
	item interface{}
	priority float64
}

type PriorityQueue struct {
	heap []priorityNode
}

func (pq PriorityQueue) IsEmpty() bool {
	return len(pq.heap) == 0
}

func (pq PriorityQueue) Size() int {
	return len(pq.heap)
}

func (pq PriorityQueue) Peek() (interface{}, float64) {
	if pq.IsEmpty() {
		panic("Cannot peek from an empty priority queue!")
	}
	return pq.heap[0].item, pq.heap[0].priority
}

func (pq *PriorityQueue) swap(i, j int) {
	temp := pq.heap[i]
	pq.heap[i] = pq.heap[j]
	pq.heap[j] = temp
}

func (pq *PriorityQueue) Insert(newItem interface{}, newPriority float64) {
	pq.heap = append(pq.heap, priorityNode{item: newItem, priority: newPriority})
	for index := len(pq.heap) - 1; index > 0; index /= 2 {
		if pq.heap[index].priority < pq.heap[index / 2].priority {
			pq.swap(index, index / 2)
		}else{
			break
		}
	}
}

func (pq *PriorityQueue) Extract() (interface{}, float64) {
	if pq.IsEmpty() {
		panic("Cannot extract from an empty priority queue!")
	}
	minNode := pq.heap[0]
	pq.swap(0, len(pq.heap) - 1)
	pq.heap = pq.heap[:(len(pq.heap) - 1)]
	for index := 0; 2 * index + 1 < len(pq.heap); {
		if 2 * index + 2 < len(pq.heap) {
			if pq.heap[2 * index + 1].priority < pq.heap[2 * index + 2].priority {
				if pq.heap[index].priority > pq.heap[2 * index + 1].priority {
					pq.swap(index, 2 * index + 1)
					index = 2 * index + 1
				}else{
					break
				}
			}else{
				if pq.heap[index].priority > pq.heap[2 * index + 2].priority {
					pq.swap(index, 2 * index + 2)
					index = 2 * index + 2
				}else{
					break
				}
			}
		}else{
			if pq.heap[index].priority > pq.heap[2 * index + 1].priority {
				pq.swap(index, 2 * index + 1)
				index = 2 * index + 1
			}else{
				break
			}
		}
	}
	return minNode.item, minNode.priority
}

func (pq *PriorityQueue) Clear() {
	pq.heap = []priorityNode{}
}