package urls

import (
    "crypto/md5"
    "encoding/json"
    "fmt"
    "github.com/CarsonSlovoka/excel/app/server"
    http2 "github.com/CarsonSlovoka/excel/pkg/net/http"
    "net/http"
)

func init() { // Initialize the directory that is specified by the user input.
    handlerFunc := func(w http.ResponseWriter, req *http.Request) {
        switch req.Method {
        case http.MethodGet:
            w.WriteHeader(http.StatusOK)
            return
        case http.MethodPost:
            maxFormSize := int64(1 << 20) // 1 MB
            if err := req.ParseMultipartForm(maxFormSize); err != nil {
                http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
                return
            }

            value := req.MultipartForm.Value
            var fileInfo FileInfo
            if val, exists := value["staticInfoObj"]; exists {
                if len(val) == 0 {
                    return
                }
                data := val[0]

                if err := json.Unmarshal([]byte(data), &fileInfo); err != nil {
                    http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
                    return
                }
            }
            staticDirID := md5.Sum([]byte(fileInfo.Path))

            prefix := fmt.Sprintf("/user/static/%x/", staticDirID)
            server.Mux.PathPrefix(prefix).Handler(
                http.StripPrefix(prefix,
                    http.FileServer(http.Dir(fileInfo.Path)),
                ),
            )
            byteData, _ := json.Marshal(map[string]interface{}{
                "staticDirURL": prefix,
            })
            _, _ = w.Write(byteData)
        }
    }
    server.Mux.HandleFunc("/user/static/", handlerFunc).Methods("GET", "POST")
}
