package urls

import (
    "fmt"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

func worksheetsHandlerFunc(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        w.WriteHeader(http.StatusOK)
        return
    }

    maxFormSize := int64(32 << 20) // 32 MB
    if err := req.ParseMultipartForm(maxFormSize); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Header().Add("Content-Type", "text/html")
        _, _ = w.Write([]byte("bad request"))
        return
    }

    fileHeader := req.MultipartForm.File["myUploadFile"][0]
    // req.MultipartForm.Value
    fmt.Println(fileHeader)
}

func init() {
    fileRouter := server.Mux.PathPrefix("/file/").Subrouter()
    fileRouter.HandleFunc("/{worksheets}", worksheetsHandlerFunc).Methods("GET", "POST")
}
