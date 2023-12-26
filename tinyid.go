package tinyid

import (
	"fmt"
)

var defaultTinyid = new(TinyId)

func SetDefault(cli *TinyId) {
	defaultTinyid = cli
}

type TinyId struct {
	services map[string]*IdService
}

func (cli *TinyId) Service(biz string) *IdService {
	if len(cli.services) > 0 {
		s, _ := cli.services[biz]
		return s
	}

	return nil
}

func (cli *TinyId) AddService(s *IdService) error {
	if s.seg == nil || !s.seg.checkInfo() {
		return fmt.Errorf("invalid id service: biz invalid")
	}

	biz := s.seg.info.Biz
	if cli.services == nil {
		cli.services = make(map[string]*IdService)
	}

	cli.services[biz] = s
	return nil
}

func Service(biz string) *IdService {
	return defaultTinyid.Service(biz)
}

func AddService(s *IdService) error {
	return defaultTinyid.AddService(s)
}
