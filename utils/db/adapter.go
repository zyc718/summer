package adapter

import (
	"encoding/json"
	"summer/utils/types"
)

// Adapter is the interface that must be implemented by a database
// adapter. The current schema supports a single connection by database type.
type Adapter interface {
	// General

	// Open and configure the
	//adapter
	Open(config json.RawMessage) error
	// TopicGet loads a single topic by name, if it exists. If the topic does not exist the call returns (nil, nil)
	TopicGet(topic string) (*types.TopicTest, error)

	// GetName returns the name of the adapter
	GetName() string
	// IsOpen checks if the adapter is ready for use
	IsOpen() bool

	// SetMaxResults configures how many results can be returned in a single DB call.
	SetMaxResults(val int) error
}
