package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

func init() {
    ht := NewTemplate("bs-table.go.html", fsTemplates, "templates/base.go.html", "templates/file/bs-table.go.html")

    ht.contextSet = append(ht.contextSet, Context{
        "Site": Setting{
            EnableBootstrap: false, // We need to use x-editable, and unfortunately, it can't support the newest version.
            EnableJquery: true, EnableFontawesome: true,
        },
        "Author":  app.Author,
        "Version": app.Version,
    })

    server.Mux.HandleFunc("/bs-table/", func(w http.ResponseWriter, r *http.Request) {
        ht.ServeHTTP(w, r)
    }).Methods("GET")
}
