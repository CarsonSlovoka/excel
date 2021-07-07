package urls

import (
    "context"
    "embed"
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

//go:embed templates/app/close-app.go.html templates/partial/jumbotron.go.html templates/base.go.html
var exitAppFS embed.FS

func init() {
    // tmpl := ParseFiles("close-app.go.html", "templates/app/close-app.go.html", "templates/partial/jumbotron.go.html", "templates/base.go.html")
    tmpl := NewTemplate("close-app.go.html", exitAppFS, "templates/app/close-app.go.html", "templates/partial/jumbotron.go.html", "templates/base.go.html")
    server.Mux.HandleFunc("/shutdown/", func(w http.ResponseWriter, r *http.Request) {
        exitHandler := func() {
            if err := server.Server.Shutdown(context.Background()); err != nil {
                log.Printf("Can't close server: %v", err)
            }
        }

        baseCtx := Context{}
        for k,v := range BaseContext {
            baseCtx[k] = v
        }
        structSite := baseCtx["Site"].(Setting)
        structSite.ShowNavbar = false
        baseCtx["Site"] = structSite

        time.AfterFunc(time.Duration(5)*time.Second, exitHandler)
        // _, _ = w.Write([]byte("Close App."))
        tmpl.contextSet = append([]Context{}, baseCtx, getLangContext(r))
        tmpl.ServeHTTP(w, r)
    })
}
