package memory

import (
	"errors"
	"github.com/uibricks/studio-engine/internal/pkg/saga/storage"
)

type memStorage struct {
	data []interface{}
}

// NewMemStorage creates log storage base on memory.
// This storage use simple `map[string][]string`, just for TestCase used.
// NOT use this in product.
func NewMemStorage() storage.Storage {
	return &memStorage{
		data: []interface{}{},
	}
}

// AppendLog appends log into queue under given logID.
func (s *memStorage) AppendLog(data interface{}) {
	s.data = append(s.data, data)
}

// Lookup return all logs
func (s *memStorage) Lookup() []interface{} {
	return s.data
}

func (s *memStorage) LastLog() (interface{}, error) {
	sizeOfLog := len(s.data)
	if sizeOfLog == 0 {
		return "", errors.New("LogData is empty")
	}
	lastLog := s.data[sizeOfLog-1]
	return lastLog, nil
}

func (s *memStorage) Cleanup() {
	s.data = make([]interface{}, 0)
}