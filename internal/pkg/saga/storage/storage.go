package storage

// Storage uses to support save and lookup saga log.
type Storage interface {

	// AppendLog appends log data into log under given logID
	AppendLog(data interface{})

	Lookup() []interface{}
	//// Lookup uses to lookup all log under given logID
	//Lookup(logID string) ([]interface{}, error)
	//
	//// Close use to close storage and release resources
	//Close() error
	//
	//// LogIDs returns exists logID
	//LogIDs() ([]string, error)
	//
	//// Cleanup cleans up all log data when end-saga is called
	Cleanup()

	// LastLog fetch last log entry with given logID
	LastLog() (interface{}, error)
}

type StorageProvider func(cfg StorageConfig) Storage

type StorageConfig struct {
}
