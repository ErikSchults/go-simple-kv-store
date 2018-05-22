package simplekvstore

import (
	"errors"
	"strings"
	"sync"
)

// ErrNotFound is returned when no key found in store
var ErrNotFound = errors.New("key does not exist")

// KeyError is returned when there is an error with a key
type KeyError struct {
	Key string
	Err error
}

func (ke KeyError) Error() string {
	return ke.Err.Error() + ": " + ke.Key
}

// KeyValue is a representation of value in the store
type KeyValue struct {
	Key   string
	Value string
}

// Store is the store
type Store struct {
	sync.Mutex
	m map[string]KeyValue
}

// New creates a new store instance with initialized storage
func New() Store {
	return Store{m: make(map[string]KeyValue)}
}

// Set sets a value to the storage while locking it
func (s *Store) Set(key, value string) {
	s.Lock()
	s.m[key] = KeyValue{Key: key, Value: value}
	s.Unlock()
}

// Get returns a value from the storage
func (s *Store) Get(key string) (KeyValue, error) {
	kv, ok := s.m[key]

	if ok {
		return kv, nil
	}

	return kv, KeyError{key, ErrNotFound}
}

// GetAll returns all kv-s from the storage
func (s *Store) GetAll() []KeyValue {
	values := make([]KeyValue, 0, len(s.m))

	for _, value := range s.m {
		values = append(values, value)
	}

	return values
}

// GetClone returns a clone of storages map
func (s *Store) GetClone() map[string]KeyValue {
	values := make(map[string]KeyValue)

	for _, value := range s.m {
		values[value.Key] = value
	}

	return values
}

// GetWithPrefix gets all values of keys starting with prefix
func (s *Store) GetWithPrefix(prefix string) []KeyValue {
	keys := []KeyValue{}

	s.Lock()
	defer s.Unlock()

	for key, kv := range s.m {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, kv)
		}
	}

	return keys
}

// GetWithPrefix gets all values of keys containing str
func (s *Store) GetKeysContaining(str string) []KeyValue {
	keys := []KeyValue{}

	s.Lock()
	defer s.Unlock()

	for key, kv := range s.m {
		if strings.Contains(key, str) {
			keys = append(keys, kv)
		}
	}

	return keys
}

// GetValuesContaining gets all values containing str
func (s *Store) GetValuesContaining(str string) []KeyValue {
	keys := []KeyValue{}

	s.Lock()
	defer s.Unlock()

	for _, kv := range s.m {
		if strings.Contains(kv.Value, str) {
			keys = append(keys, kv)
		}
	}

	return keys
}
