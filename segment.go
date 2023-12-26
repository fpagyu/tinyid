package tinyid

import "sync/atomic"

type Segment struct {
	info   *IdInfo
	min_id int64
	max_id int64
	src    IdSource

	sem   int32
	cache *Segment
}

func (seg *Segment) Id() int64 {
	id := atomic.AddInt64(&seg.min_id, 1)
	if id <= seg.max_id {
		if seg.cache == nil {
			seg.updateCache()
		}

		if bits := seg.info.Pid_Bits; bits > 0 {
			id <<= int64(bits)
			id |= int64(seg.info.Pid)
		}
		return id
	}

	return 0
}

func (seg *Segment) checkInfo() bool {
	return seg.info != nil && seg.info.Biz != ""
}

func (seg *Segment) setIdInfo(info *IdInfo) {
	if info != nil {
		seg.info = info
		seg.max_id = info.MaxId
		seg.min_id = info.MinId()
	}
}

// 更新号段
func (seg *Segment) updateSegment() {
	if seg.cache != nil {
		seg.setIdInfo(seg.cache.info)
		seg.cache = nil
		return
	}

	info := seg.src.Get(seg.info.Biz)
	if info != nil {
		seg.setIdInfo(info)
	}
}

func (seg *Segment) valid() bool {
	return seg.min_id < seg.max_id
}

func (seg *Segment) usage() float64 {
	d := float64(seg.max_id - seg.min_id)
	return 1 - d/float64(seg.info.Step)
}

func (seg *Segment) updateCache() {
	// 当前号段使用率超过70%, 异步更新缓存号段
	if seg.usage() < 0.7 {
		return
	}

	if n := atomic.AddInt32(&seg.sem, 1); n == 1 {
		// update cache segment
		go func() {
			defer atomic.AddInt32(&seg.sem, -1)

			info := seg.src.Get(seg.info.Biz)
			seg.cache = &Segment{info: info}
		}()
	} else {
		atomic.AddInt32(&seg.sem, -1)
	}
}
