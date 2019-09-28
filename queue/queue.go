package queue

type Queue struct {
	nodes [][]byte
	size  int
	head  int
}

func (q *Queue) Push(n []byte) {
	if q.head >= q.size-1 {
		for i := 0; i < q.size-1; i++ {
			q.nodes[i] = q.nodes[i+1]
		}
		q.head = q.size - 1
	}
	q.nodes[q.head] = n
	q.head++
}

func (q *Queue) GetAll() [][]byte {
	return q.nodes
}

func (q *Queue) GetSize() int {
	return q.size
}

func (q *Queue) GetLength() int {
	return q.head
}

func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([][]byte, size),
		size:  size,
	}
}
