package tinyid

import (
	"database/sql"
	"log"
)

type IdSource interface {
	Get(biz string) *IdInfo
}

type SqlDB struct {
	*sql.DB
}

func NewSqlDB(db *sql.DB) IdSource {
	return &SqlDB{DB: db}
}

func (dbs *SqlDB) Get(biz string) *IdInfo {
	var info IdInfo
	row := dbs.QueryRow("SELECT `id`, `max_id`, `step`, `version`, `pid`, `pid_bits` FROM `id_info` WHERE `biz` = ?", biz)
	if err := row.Scan(&info.Id, &info.MaxId, &info.Step, &info.Version, &info.Pid, &info.Pid_Bits); err != nil {
		log.Println("get id_info:", err)
		return nil
	}

	version := info.Version + 1
	max_id := info.MaxId + int64(info.Step)
	r, err := dbs.Exec("UPDATE `id_info` SET `max_id`=?, `version`=? WHERE `biz` = ? AND `version` = ?",
		max_id, version, biz, info.Version)
	if err != nil {
		log.Println("update id_info:", err)
		return nil
	}

	if n, _ := r.RowsAffected(); n <= 0 {
		log.Println("update id_info: no id_info row affect")
		return nil
	}

	info.Biz = biz
	info.MaxId = max_id
	return &info
}
