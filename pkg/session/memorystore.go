package session

import "time"

type MemoryStore struct {
    provider     Provider
    id           string
    value        map[interface{}]interface{}
    timeAccessed time.Time
}

func (s *MemoryStore) Set(key, value interface{}) error {
    s.value[key] = value
    return s.provider.SessionUpdate(s.id)
}

func (s *MemoryStore) Get(key interface{}) interface{} {
    if err := s.provider.SessionUpdate(s.id); err != nil {
        panic(err)
    }

    if v, ok := s.value[key]; ok {
        return v
    } else {
        return nil
    }
}

func (s *MemoryStore) GetMap() map[interface{}]interface{} {
    if err := s.provider.SessionUpdate(s.id); err != nil {
        panic(err)
    }
    return s.value
}

func (s *MemoryStore) Clear(key interface{}) error {
    delete(s.value, key)
    return s.provider.SessionUpdate(s.id)
}

func (s *MemoryStore) SessionID() string {
    return s.id
}
