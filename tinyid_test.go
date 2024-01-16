package tinyid

import (
	"database/sql"
	"log"
	"sync"
	"testing"
	// _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func setup() {
	var err error
	dataSourceName := "xxxxxx"
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("connect db success-g!!!")
}

func teardown() {
	if db != nil {
		db.Close()
	}
}

func Test_NextId(t *testing.T) {
	setup()
	defer teardown()

	// id := NewIdService("pd_asset", NewSqlDB(db))
	// id := (IdInfo{}).NewIdService(nil)
	// AddService(IdInfo{
	// 	Biz: "pd_asset",
	// }.NewIdService(NewSqlDB(db)))
	Init(db)
	AddService("pd_asset")

	id := Service("pd_asset")
	var wg sync.WaitGroup
	for j := 0; j < 10; j++ {
		gid := j
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 90; i++ {
				// ids := id.NextIds(10)
				// t.Log("[", ids[0], ",", ids[9], "]")
				t.Logf("g%d: %d", gid, id.NextId())
			}
		}()
	}
	wg.Wait()

	// for i := 0; i < 300; i++ {
	// 	t.Logf("id-%d: %d", i, id.NextId())
	// }
}
