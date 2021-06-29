package server

import (
    "fmt"
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/pkg/session"
    "github.com/gorilla/mux"
    "net/http"
)

var (
    Mux           *mux.Router // http.HandleFunc
    Server        http.Server
    SessionManger *session.Manager
)

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

func ListenAndServe() error {
    Server = http.Server{Addr: fmt.Sprintf(":%s", app.Port), Handler: Mux}
    err := Server.ListenAndServe()
    return err
}
