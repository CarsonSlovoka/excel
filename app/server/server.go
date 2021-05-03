package server

import (
    "github.com/CarsonSlovoka/excel/api/session"
    "github.com/gorilla/mux"
    "net/http"
)

var (
    Mux          *mux.Router // http.HandleFunc
    server        http.Server
    SessionManger *session.Manager
)

func GetServer() *http.Server {
    return &server
}

func init() {
    Mux = mux.NewRouter()
}

func init() {
    var err error
    providerName := "my-memory-provider"
    cookieName := "session"
    session.Register(providerName, session.GetMemoryProvider())
    SessionManger, err = session.NewManager(providerName, cookieName, 300) // 5 minute
    if err != nil {
        panic(err)
    }
    go SessionManger.GC() // start GC Processing
}
