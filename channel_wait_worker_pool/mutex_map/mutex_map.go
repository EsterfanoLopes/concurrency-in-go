package mutex_map

import "sync"

// MutexMap is a map with a mutex. It is safe for concurrent use
type MutexMap struct {
	rwlock *sync.RWMutex
	m      map[uint]interface{} // map to store the data
}

// New creates a new MutexMap
func New() *MutexMap {
	return &MutexMap{
		rwlock: &sync.RWMutex{},
		m:      make(map[uint]interface{}),
	}
}

// Get returns the value and if the key exists
func (mm *MutexMap) Get(key uint) (interface{}, bool) {
	mm.rwlock.RLock()
	defer mm.rwlock.RUnlock()
	value, ok := mm.m[key]
	return value, ok
}

// Set sets the value for the key
func (mm *MutexMap) Set(key uint, value interface{}) {
	mm.rwlock.Lock()
	defer mm.rwlock.Unlock()
	mm.m[key] = value
}

// Delete deletes the key
func (mm *MutexMap) Delete(key uint) {
	mm.rwlock.Lock()
	defer mm.rwlock.Unlock()
	delete(mm.m, key)
}
