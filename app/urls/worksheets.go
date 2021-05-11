package urls

import (
    "encoding/json"
    "github.com/CarsonSlovoka/excel/app/server"
    fileHelper "github.com/CarsonSlovoka/excel/pkg/file"
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

    server.Mux.HandleFunc(req.URL.Path+"/data", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")

        file, err := fileHeader.Open()
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            _, _ = w.Write([]byte("Couldn't open the file"))
            return
        }
        defer func() { _ = file.Close() }()
        // byteData, err := ioutil.ReadAll(file)

        dataArray, err := fileHelper.CSV2Json(file)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            _, _ = w.Write([]byte(err.Error()))
            return
        }
        // beautifulJsonByte, err := json.MarshalIndent(dataArray, "", "  ")
        jsonByte, err := json.Marshal(dataArray)

        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            _, _ = w.Write([]byte(err.Error()))
        }
        _, _ = w.Write(jsonByte)
    })
}

func init() {
    fileRouter := server.Mux.PathPrefix("/file/").Subrouter()
    fileRouter.HandleFunc("/{worksheets}", worksheetsHandlerFunc).Methods("GET", "POST")
}
