package urls

import (
    "embed"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

const (
    sessionName = "session"
)

//go:embed templates/base.go.html templates/config/popupConfig.go.html
//go:embed templates/file/*.go.html
var fsTemplates embed.FS

func initFileURL() {
    ht := NewTemplate("file.go.html", fsTemplates, "templates/base.go.html", "templates/file/file.go.html")
    server.Mux.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
        ht.contextSet = append([]Context{}, BaseContext, getLangContext(r), Context{
            "Params": struct {
                CSSList []string
                JSList  []string
            }{},
        })

        ht.ServeHTTP(w, r)
    },
    ).Methods("GET")
}
