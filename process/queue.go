package process

import "simulator/models"

type Queue struct {
	items []*models.Process
}

func (q *Queue) Enqueue(p *models.Process) { q.items = append(q.items, p) }

func (q *Queue) Dequeue() *models.Process {
	if len(q.items) == 0 {
		return nil
	}
	p := q.items[0]
	q.items = q.items[1:]
	return p
}

func (q *Queue) Len() int { return len(q.items) }

func (q *Queue) Snapshot() []*models.Process {
	cp := make([]*models.Process, len(q.items))
	copy(cp, q.items)
	return cp
}
