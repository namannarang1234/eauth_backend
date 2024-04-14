package types

import "sync"

type SafeMap struct {
	mt sync.RWMutex
	d  map[string]bool
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		d:  make(map[string]bool),
		mt: sync.RWMutex{},
	}
}

func (m *SafeMap) Put(key string, value bool) {
	m.mt.Lock()
	defer m.mt.Unlock()

	m.d[key] = true
}

func (m *SafeMap) Get(key string) (bool, bool) {
	m.mt.RLock()
	defer m.mt.RUnlock()

	if v, ok := m.d[key]; ok {
		return v, true
	} else {
		return false, false
	}
}

func (m *SafeMap) Delete(key string) {
	m.mt.Lock()
	defer m.mt.Unlock()

	delete(m.d, key)
}
