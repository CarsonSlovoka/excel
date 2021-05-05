package urls

import (
    "encoding/json"
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
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
