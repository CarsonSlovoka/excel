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

func (mP *MemoryProvider) SessionInit(id string) (Session, error) {
    mP.lock.Lock()
    defer mP.lock.Unlock()
    newSessStore := &Store{
        provider:     mP,
        id:           id,
        timeAccessed: time.Now(),
        value:        make(map[interface{}]interface{}, 0),
    }
    element := mP.list.PushBack(newSessStore)
    mP.mapStore[id] = element
    return newSessStore, nil
}

func (mP *MemoryProvider) SessionRead(id string) (Session, error) {
    if element, ok := mP.mapStore[id]; ok {
        return element.Value.(*Store), nil
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
        if (element.Value.(*Store).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
            mP.list.Remove(element)
            delete(mP.mapStore, element.Value.(*Store).id)
        } else {
            break
        }
    }
}

func (mP *MemoryProvider) SessionUpdate(id string) error {
    mP.lock.Lock()
    defer mP.lock.Unlock()
    if element, ok := mP.mapStore[id]; ok {
        element.Value.(*Store).
            timeAccessed = time.Now()
        mP.list.MoveToFront(element)
        return nil
    }
    return nil
}
