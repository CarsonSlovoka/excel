package urls

import (
    "embed"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

const (
    sessionName = "session"
)

//go:embed templates/base.go.html
//go:embed templates/file/*.go.html
var fsTemplates embed.FS

//go:embed static/css/file/*.css
var fsLoginCSS embed.FS

//go:embed static/js/file/*.js
var fsLoginJS embed.FS

func initFileURL() {
    ht := NewTemplate("file.go.html", fsTemplates, "templates/base.go.html", "templates/file/file.go.html")
    ht.contextSet = append(ht.contextSet, BaseContext, Context{
        "Params": struct {
            CSSList []string
            JSList  []string
        }{},
    })
    server.Mux.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
        ht.ServeHTTP(w, r)
    },
    ).Methods("GET")
}
