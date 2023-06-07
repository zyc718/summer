package store

import (
	"encoding/json"
	"errors"
	"fmt"
	adapter "summer/utils/db"
	"summer/utils/types"
)

type configType struct {
	// 16-byte key for XTEA. Used to initialize types.UidGenerator.
	UidKey []byte `json:"uid_key"`
	// Maximum number of results to return from adapter.
	MaxResults int `json:"max_results"`
	// DB adapter name to use. Should be one of those specified in `Adapters`.
	UseAdapter string `json:"use_adapter"`
	// Configurations for individual adapters.
	Adapters map[string]json.RawMessage `json:"adapters"`
}

var adp adapter.Adapter

var availableAdapters = make(map[string]adapter.Adapter)

func OpenAdapter(workerId int, jsonconf json.RawMessage) error {
	var config configType
	err := json.Unmarshal(jsonconf, &config)

	if err != nil {
		return err
	}
	if adp == nil {
		if len(config.UseAdapter) > 0 {
			if ad, ok := availableAdapters[config.UseAdapter]; ok {
				adp = ad
			} else {
				return errors.New("store: " + config.UseAdapter + " adapter is not available in this binary")
			}

		} else if len(availableAdapters) == 1 {
			fmt.Printf("availableAdapters is %v\n", availableAdapters)
			// Default to the only entry in availableAdapters.
			for _, v := range availableAdapters {
				adp = v
			}
		} else {
			fmt.Printf("store OpenAdapter is error  \n")
			return errors.New("store: db adapter is not specified. Please set `store_config.use_adapter` in `tinode.conf`")
		}
	}
	//fmt.Printf("结果是%+v\n", config)
	if adp.IsOpen() {
		return errors.New("store: connection is already opened")
	}
	// Initialise snowflake
	if workerId < 0 || workerId > 1023 {
		return errors.New("store: invalid worker ID")
	}

	if err := adp.SetMaxResults(config.MaxResults); err != nil {
		return err
	}

	var adapterConfig json.RawMessage
	if config.Adapters != nil {
		adapterConfig = config.Adapters[adp.GetName()]
	}
	return adp.Open(adapterConfig)
}

type TopicsObjMapper struct{}

var Topics TopicsObjMapper

func (TopicsObjMapper) Get(topic string) (*types.Topic, error) {
	res, err := adp.TopicGet(topic)
	return res, err
}

// GetAdapterName returns the name of the current adater.
func GetAdapterName() string {
	if adp != nil {
		return adp.GetName()
	}

	return ""
}

// RegisterAdapter makes a persistence adapter available.
// If Register is called twice or if the adapter is nil, it panics.
func RegisterAdapter(a adapter.Adapter) {
	if a == nil {
		panic(any("store: Register adapter is nil"))
	}
	adapterName := a.GetName()
	if _, ok := availableAdapters[adapterName]; ok {
		panic(any("store: adapter '" + adapterName + "' is already registered"))
	}
	availableAdapters[adapterName] = a
}
