package tinyid

type IdInfo struct {
	Id  int
	Biz string

	MaxId   int64
	Step    int // 每次获取号段的长度
	Version int // 乐观锁

	Pid      int  // 分片id
	Pid_Bits int8 // 分片id比特数
}

func (info *IdInfo) MinId() int64 {
	v := info.MaxId - int64(info.Step)
	if v < 0 {
		return 0
	} else {
		return v
	}
}

func (info IdInfo) NewIdService(src IdSource) *IdService {
	segment := &Segment{
		src: src,
	}
	segment.setIdInfo(&info)

	return &IdService{seg: segment}
}
