package queue_range

type QueueRange struct {
	ranges []Range
}

type Range struct {
	start int
	end   int
}
