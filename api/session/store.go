package session

import "time"

type Store struct {
    provider     Provider
    id           string
    value        map[interface{}]interface{}
    timeAccessed time.Time
}

func (s *Store) Set(key, value interface{}) error {
    s.value[key] = value
    return s.provider.SessionUpdate(s.id)
}

func (s *Store) Get(key interface{}) interface{} {
    if err := s.provider.SessionUpdate(s.id); err != nil {
        panic(err)
    }

    if v, ok := s.value[key]; ok {
        return v
    } else {
        return nil
    }
}

func (s *Store) Clear(key interface{}) error {
    delete(s.value, key)
    return s.provider.SessionUpdate(s.id)
}

func (s *Store) SessionID() string {
    return s.id
}
