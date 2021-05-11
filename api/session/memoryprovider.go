package session

import (
    "container/list"
    "errors"
    "sync"
    "time"
)

var (
    memoryProvider *MemoryProvider
)

func init() {
    memoryProvider = &MemoryProvider{
        mapStore: make(map[string]*list.Element, 1),
        list:     list.New(),
    }
}

func GetMemoryProvider() *MemoryProvider {
    return memoryProvider
}

var (
    SIDNotFoundError error
)

func init() {
    SIDNotFoundError = errors.New("session id not exits")
}

type MemoryProvider struct {
    mapStore map[string]*list.Element
    list     *list.List
    lock     sync.Mutex
}

func (mP *MemoryProvider) SessionInit(id string) (Store, error) {
    mP.lock.Lock()
    defer mP.lock.Unlock()
    memoryStore := &MemoryStore{
        provider:     mP,
        id:           id,
        timeAccessed: time.Now(),
        value:        make(map[interface{}]interface{}, 0),
    }
    element := mP.list.PushBack(memoryStore)
    mP.mapStore[id] = element
    return memoryStore, nil
}

func (mP *MemoryProvider) SessionRead(id string) (Store, error) {
    if element, ok := mP.mapStore[id]; ok {
        return element.Value.(*MemoryStore), nil
    } else {
        return nil, SIDNotFoundError
    }
}

func (mP *MemoryProvider) SessionDestroy(id string) error {
    if element, ok := mP.mapStore[id]; ok {
        delete(mP.mapStore, id)
        mP.list.Remove(element)
        return nil
    }
    return nil
}

func (mP *MemoryProvider) SessionGC(maxLifeTime int64) {
    mP.lock.Lock()
    defer mP.lock.Unlock()

    for {
        element := mP.list.Back()
        if element == nil {
            break
        }
        if (element.Value.(*MemoryStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
            mP.list.Remove(element)
            delete(mP.mapStore, element.Value.(*MemoryStore).id)
        } else {
            break
        }
    }
}

func (mP *MemoryProvider) SessionUpdate(id string) error {
    mP.lock.Lock()
    defer mP.lock.Unlock()
    if element, ok := mP.mapStore[id]; ok {
        element.Value.(*MemoryStore).
            timeAccessed = time.Now()
        mP.list.MoveToFront(element)
        return nil
    }
    return nil
}
