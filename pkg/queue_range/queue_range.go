package queue_range

import (
	"fmt"
	"slices"
)

func NewQueueRange() *QueueRange {
	return &QueueRange{
		ranges: nil,
	}
}

func (q *QueueRange) AddToQueue(id int) error {
	if id < 0 {
		return fmt.Errorf("error adding id of %v must be non-negative", id)
	}

	if q.ranges == nil {
		q.ranges = make([]Range, 1)
		q.ranges[0].start = id
		q.ranges[0].end = id
		return nil
	}

	rangeToEdit := 0
	for i := len(q.ranges) - 1; i > -1; i-- {
		if id > q.ranges[i].start {
			rangeToEdit = i
			break
		}
	}

	if q.ranges[rangeToEdit].end+1 == id {
		q.ranges[rangeToEdit].end++
	} else if q.ranges[rangeToEdit].end >= id {
		return nil
	} else if id > q.ranges[rangeToEdit].end {
		newRange := Range{
			start: id,
			end:   id,
		}
		q.ranges = slices.Insert(q.ranges, rangeToEdit+1, newRange)
	}

	return nil
}

func (q *QueueRange) GetFirstInQueue() int {
	if q.ranges == nil {
		return 0
	}
	return q.ranges[0].start
}

func (q *QueueRange) GetLastInQueue() int {
	if q.ranges == nil {
		return 0
	}
	return q.ranges[len(q.ranges)-1].end
}

func (q *QueueRange) GetNextInQueue(prev int) (int, error) {
	if q.ranges == nil {
		return 0, fmt.Errorf("error: queue is empty when trying to get next for prev %v", prev)
	} else if prev < q.ranges[0].start {
		return q.ranges[0].start, nil
	}

	next := -1
	for i := 0; i < len(q.ranges); i++ {
		if prev >= q.ranges[i].start && prev < q.ranges[i].end {
			next = prev + 1
			break
		}
	}
	if next == -1 {
		for i := 0; i < len(q.ranges); i++ {
			if prev == q.ranges[i].end {
				if len(q.ranges) == i+1 {
					return -1, fmt.Errorf("error: the end of the queue was reached at prev: %v", prev)
				} else {
					next = q.ranges[i+1].start
					break
				}
			}
		}
	}
	return next, nil
}

func (q *QueueRange) TotalItemsInQueue() int {
	if q.ranges == nil {
		return 0
	}
	sum := 0
	for i := 0; i < len(q.ranges); i++ {
		sum += q.ranges[i].end - q.ranges[i].start + 1
	}
	return sum
}
