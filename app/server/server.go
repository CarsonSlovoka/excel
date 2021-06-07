package server

import (
    "context"
    "fmt"
    "github.com/CarsonSlovoka/excel/pkg/session"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "time"
)

var (
    Mux           *mux.Router // http.HandleFunc
    server        http.Server
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
    Mux.HandleFunc("/shutdown/", func(w http.ResponseWriter, r *http.Request) {
        exitHandler := func() {
            if err := server.Shutdown(context.Background()); err != nil {
                log.Printf("Can't close server: %v", err)
            }
        }

        time.AfterFunc(time.Duration(5)*time.Second, exitHandler)
        _, _ = w.Write([]byte("Close App."))
    })

    server = http.Server{Addr: fmt.Sprintf(":%s", "7121"), Handler: Mux}
    err := server.ListenAndServe()
    return err
}
