package types

import "time"

// TimeNow returns current wall time in UTC rounded to milliseconds.
//返回毫秒
func TimeNow() time.Time {
	//	设置时区
	l, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().UTC().Round(time.Millisecond).In(l)
}

type Uid uint64
type ObjState int
type AccessMode uint

// StringSlice is defined so Scanner and Valuer can be attached to it.
type StringSlice []string

// ObjHeader is the header shared by all stored objects.
type ObjHeader struct {
	// using string to get around rethinkdb's problems with uint64;
	// `bson:"_id"` tag is for mongodb to use as primary key '_id'.
	Id        string `bson:"_id" db:"id"`
	id        Uid
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultAccess is a per-topic default access modes
type DefaultAccess struct {
	Auth AccessMode
	Anon AccessMode
}

type perUserData struct {
	private interface{}
	want    AccessMode
	given   AccessMode
}

// Topic stored in database. Topic's name is Id
type Topic struct {
	ObjHeader `bson:",inline"`

	// State of the topic: normal (ok), suspended, deleted
	State   ObjState
	StateAt *time.Time `json:"StateAt,omitempty" bson:",omitempty"`

	// Timestamp when the last message has passed through the topic
	TouchedAt time.Time

	// Use bearer token or use ACL
	UseBt bool

	// Topic owner. Could be zero
	Owner string

	// Default access to topic
	Access DefaultAccess

	// Server-issued sequential ID
	SeqId int
	// If messages were deleted, sequential id of the last operation to delete them
	DelId int

	Public interface{}

	// Indexed tags for finding this topic.
	Tags StringSlice

	// Deserialized ephemeral params
	perUser map[Uid]*perUserData // deserialized from Subscription
}
