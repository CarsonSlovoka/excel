package urls

import (
    "context"
    "encoding/json"
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "log"
    "net/http"
    "time"
)

func initSystemURL() {
    serverMux := server.Mux

    serverMux.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
        beautifulJsonByte, err := json.MarshalIndent(app.About, "", "  ")
        if err != nil {
            panic(err)
        }
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write(beautifulJsonByte)
    }).Methods("GET")
}

func init() {
    server.Mux.HandleFunc("/shutdown/", func(w http.ResponseWriter, r *http.Request) {
        exitHandler := func() {
            if err := server.Server.Shutdown(context.Background()); err != nil {
                log.Printf("Can't close server: %v", err)
            }
        }

        time.AfterFunc(time.Duration(5)*time.Second, exitHandler)
        _, _ = w.Write([]byte("Close App."))
    })
}
