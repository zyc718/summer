package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"summer/utils/db/pool"
	"summer/utils/store"
	"summer/utils/types"
	"time"
)

// adapter holds MySQL connection data.
type adapter struct {
	db     *sqlx.DB
	dsn    string
	dbName string
	// Maximum number of records to return
	maxResults int
	// Maximum number of message records to return
	maxMessageResults int
	version           int
}

const (
	defaultDSN      = "root:@tcp(localhost:3306)/tinode?parseTime=true"
	defaultDatabase = "tinode"

	adpVersion = 111

	adapterName = "local"

	defaultMaxResults = 1024
	// This is capped by the Session's send queue limit (128).
	defaultMaxMessageResults = 100
)

func (a *adapter) TopicGet(topic string) (*types.Topic, error) {
	//TODO implement me
	fmt.Printf("这是mysql TopicGet \n")
	type Test struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	var dest Test
	err := a.db.QueryRowx("select id ,name from test where name = ?", topic).StructScan(&dest)
	if err != nil {
		return nil, err
	}
	return &types.Topic{
		Owner: dest.Name,
	}, nil
}

func (a *adapter) GetName() string {
	//TODO implement me
	fmt.Printf("这是mysql GetName \n")
	return adapterName

}

func (a *adapter) Open(config json.RawMessage) error {
	fmt.Printf("这是 mysql 包的open 方法 ")
	if a.maxResults <= 0 {
		a.maxResults = defaultMaxResults
	}

	if a.maxMessageResults <= 0 {
		a.maxMessageResults = defaultMaxMessageResults
	}

	db, err := pool.Pool.Open(adapterName)

	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(120 * time.Second)

	a.db = db

	return nil
}

// IsOpen returns true if connection to database has been established. It does not check if
// connection is actually live.
func (a *adapter) IsOpen() bool {
	return a.db != nil
}

// SetMaxResults configures how many results can be returned in a single DB call.
func (a *adapter) SetMaxResults(val int) error {
	if val <= 0 {
		a.maxResults = defaultMaxResults
	} else {
		a.maxResults = val
	}

	return nil
}

func init() {
	fmt.Printf("这是mysql 的init 方法 \n")
	store.RegisterAdapter(&adapter{})
}
