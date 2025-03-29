package queue_range

import (
	"fmt"
	"testing"
)

func TestQueueRange(t *testing.T) {
	q := NewQueueRange()

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	if q.AddToQueue(5) != nil {
		t.Errorf("value 5 not added to queue correctly")
	}

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	for i := 5; i < 20; i++ {
		if q.AddToQueue(i) != nil {
			t.Errorf("value 6 not added to queue correctly")
		}
	}

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	for i := 30; i < 50; i++ {
		if q.AddToQueue(i) != nil {
			t.Errorf("value 6 not added to queue correctly")
		}
	}

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	for i := 25; i < 27; i++ {
		if q.AddToQueue(i) != nil {
			t.Errorf("value 6 not added to queue correctly")
		}
	}

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	if q.AddToQueue(27) != nil {
		t.Errorf("value 27 not added to queue correctly")
	}
	if q.AddToQueue(24) != nil {
		t.Errorf("value 24 not added to queue correctly")
	}
	if q.AddToQueue(22) != nil {
		t.Errorf("value 24 not added to queue correctly")
	}

	fmt.Printf("Start of queue: %v End of queue: %v Size of queue: %v\n", q.GetFirstInQueue(), q.GetLastInQueue(), q.TotalItemsInQueue())

	for i := 0; i < len(q.ranges); i++ {
		fmt.Printf("Range start %v range end %v\n", q.ranges[i].start, q.ranges[i].end)
	}
	fmt.Printf("Number of ranges: %v\n", len(q.ranges))

	val1, err1 := q.GetNextInQueue(0)
	if err1 != nil {
		t.Errorf("next value 0 not given correctly")
	}
	fmt.Printf("Get next in queue 0: %v\n", val1)
	val2, err2 := q.GetNextInQueue(22)
	if err2 != nil {
		t.Errorf("next value 22 not given correctly")
	}
	fmt.Printf("Get next in queue 22: %v\n", val2)
	val3, err3 := q.GetNextInQueue(24)
	if err3 != nil {
		t.Errorf("next value 24 not given correctly")
	}
	fmt.Printf("Get next in queue 24: %v\n", val3)
	val4, err4 := q.GetNextInQueue(48)
	if err4 != nil {
		t.Errorf("next value 48 not given correctly")
	}
	fmt.Printf("Get next in queue 48: %v\n", val4)
}
