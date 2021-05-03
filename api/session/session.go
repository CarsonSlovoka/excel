package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Clear(key interface{}) error
	SessionID() string
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionUpdate(sid string) error // for fast GC and query
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

var mapProvider map[string]Provider

func init() {
	mapProvider = make(map[string]Provider)
}

func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := mapProvider[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	mapProvider[name] = provide
}


type Manager struct {
	provider    Provider
	cookieName  string
	maxLifeTime int64
	lock        sync.Mutex
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := mapProvider[provideName]
	if !ok {
        // %q: a single-quoted character literal safely escaped with Go syntax.
		return nil, fmt.Errorf(`session: unknown "Provider" please Regester before calling this function. %q`, provideName) // https://golang.org/pkg/fmt/
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err == nil && cookie.Value != "" {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err = manager.provider.SessionRead(sid)
		if err == nil {
			return session
		}
	}
	sid := manager.generateSessionID()
	session, err = manager.provider.SessionInit(sid)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	cookie = &http.Cookie{
		Name:     manager.cookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(manager.maxLifeTime)}
	http.SetCookie(w, cookie)

	return session
}

func (manager *Manager) SessionQuery(r *http.Request) (session Session) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		return nil
	}
	sid, _ := url.QueryUnescape(cookie.Value)
	session, _ = manager.provider.SessionRead(sid)
	return session
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		_ = manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(
		time.Duration(manager.maxLifeTime)*time.Second,
		func() { manager.GC() },
	)
}

func (manager *Manager) generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
