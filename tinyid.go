package tinyid

import (
	"database/sql"
	"fmt"
	"sync"
)

var defaultTinyid = new(TinyId)

func SetDefault(cli *TinyId) {
	defaultTinyid = cli
}

type TinyId struct {
	idsource IdSource
	services sync.Map
}

func (cli *TinyId) SetIdSource(src IdSource) {
	cli.idsource = src
}

func (cli *TinyId) Service(biz string) *IdService {
	if s, ok := cli.services.Load(biz); ok {
		return s.(*IdService)
	}

	return nil
}

func (cli *TinyId) AddService(biz string) (*IdService, error) {
	if biz == "" {
		return nil, fmt.Errorf("biz cannot be empty")
	}

	if cli.idsource == nil {
		return nil, fmt.Errorf("tinyid is not initialized")
	}

	info := IdInfo{Biz: biz}
	service := NewIdService(&info, cli.idsource)
	s, _ := cli.services.LoadOrStore(biz, service)

	return s.(*IdService), nil
}

func Init(db *sql.DB) {
	defaultTinyid.SetIdSource(NewSqlDB(db))
}

func Service(biz string) *IdService {
	return defaultTinyid.Service(biz)
}

func AddService(biz string) (*IdService, error) {
	return defaultTinyid.AddService(biz)
}
