package tinyid

import (
	"sync"
)

type IdService struct {
	seg *Segment

	lock sync.RWMutex
}

func (id *IdService) updateSegment() {
	id.lock.Lock()
	defer id.lock.Unlock()

	if id.seg.valid() {
		return
	}

	id.seg.updateSegment()
}

func (id *IdService) nextId() (r int64) {
	id.lock.RLock()
	r = id.seg.Id()
	id.lock.RUnlock()

	return
}

func (id *IdService) NextId() int64 {
	n := 5
	for n > 0 {
		if n := id.nextId(); n > 0 {
			return n
		}

		n--
		if !id.seg.valid() {
			id.updateSegment()
		}
	}

	return 0
}

func (id *IdService) NextIds(n int) []int64 {
	ids := make([]int64, n)

	for i := range ids {
		ids[i] = id.NextId()
	}

	return ids
}
