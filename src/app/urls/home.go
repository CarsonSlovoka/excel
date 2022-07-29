package urls

import (
	"github.com/CarsonSlovoka/excel/app/server"
	"net/http"
)

func initHomeURL() {
	server.Mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/file/", http.StatusSeeOther) // It's allow you to open the file (CSV format)
		})
}
