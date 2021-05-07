package urls

import (
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)


func init() {
    ht := NewTemplate("bs-table.go.html", fsTemplates, "templates/base.go.html", "templates/file/bs-table.go.html")
    ht.contextSet = append(ht.contextSet, BaseContext)

    server.Mux.HandleFunc("/bs-table/", func(w http.ResponseWriter, r *http.Request) {
        ht.ServeHTTP(w, r)
    }).Methods("GET")
}
