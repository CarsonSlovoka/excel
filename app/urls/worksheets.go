package urls

import (
    "fmt"
    "github.com/CarsonSlovoka/excel/app/server"
    http2 "github.com/CarsonSlovoka/excel/pkg/net/http"
    "net/http"
)

func worksheetsHandlerFunc(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        w.WriteHeader(http.StatusOK)
        return
    }

    maxFormSize := int64(32 << 20) // 32 MB
    if err := req.ParseMultipartForm(maxFormSize); err != nil {
        http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
        return
    }

    // fileHeader := req.MultipartForm.File["myUploadFile"][0]
    value := req.MultipartForm.Value
    if val, exists := value["uploadData"]; exists {
        if len(val) == 0 {
            return
        }
        data := val[0]
        fmt.Println(data)
    }
}

func init() {
    fileRouter := server.Mux.PathPrefix("/file/").Subrouter()
    fileRouter.HandleFunc("/{worksheets}", worksheetsHandlerFunc).Methods("GET", "POST")
}
