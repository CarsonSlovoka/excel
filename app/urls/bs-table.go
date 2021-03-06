package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

func init() {
    // {{template "popupConfig.go.html" .}}
    ht := NewTemplate("bs-table.go.html", fsTemplates, "templates/base.go.html", "templates/file/bs-table.go.html", "templates/config/popupConfig.go.html")

    server.Mux.HandleFunc("/bs-table/", func(w http.ResponseWriter, r *http.Request) {
        ht.contextSet = append([]Context{}, getLangContext(r), Context{
            "Site": Setting{
                EnableBootstrap: false, // We need to use x-editable, and unfortunately, it can't support the newest version.
                EnableJquery:    true, EnableFontawesome: true,
                ShowFooter: false,
            },
            "Author":  app.Author,
            "Version": app.Version,
        })
        ht.ServeHTTP(w, r)
    }).Methods("GET")
}
