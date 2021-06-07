package urls

import (
    "encoding/json"
    "errors"
    "github.com/CarsonSlovoka/excel/app/server"
    http2 "github.com/CarsonSlovoka/excel/pkg/net/http"
    "net/http"
    "os"
    "strings"
    "time"
)

type FileInfo struct {
    Name    string    // base name of the file
    Size    int64     // length in bytes for regular files; system-dependent for others
    ModTime time.Time // modification time
    IsDir   bool      // abbreviation for Mode().IsDir()
    Path    string    // absolute path
}

func funcHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        w.WriteHeader(http.StatusOK)
        return
    }

    maxFormSize := int64(32 << 20) // 32 MB
    if err := req.ParseMultipartForm(maxFormSize); err != nil {
        http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
        return
    }

    urlPath := req.URL.Path
    funcName := strings.TrimLeft(urlPath, "/api/")
    var mapVal map[string][]string
    mapVal = req.MultipartForm.Value

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    switch funcName {
    case "os/Stat":
        var inputPara string
        for key, _ := range mapVal {
            inputPara = key
            break
        }

        if val, err := getStringFromValue(mapVal, inputPara); err == nil {
            inputPath := val[0]
            fileInfo, err := os.Stat(inputPath)
            if err != nil {
                http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
                return
            }
            result := &FileInfo{
                fileInfo.Name(),
                fileInfo.Size(),
                fileInfo.ModTime(),
                fileInfo.IsDir(),
                inputPath,
            }
            byteData, _ := json.Marshal(result)
            _, _ = w.Write(byteData)
        }
    default:
        http2.ShowErrorRequest(w, http.StatusBadRequest, "wrong function name!")
    }
}

func getStringFromValue(valMap map[string][]string, targetName string) ([]string, error) {
    if val, exists := valMap[targetName]; exists {
        if len(val) != 0 {
            return val, nil
        }
    }
    return nil, errors.New("value not found error")
}

func init() {
    apiOSRouter := server.Mux.PathPrefix("/api/os/").Subrouter()
    apiOSRouter.HandleFunc("/{funName}", funcHandler).Methods("GET", "POST")
}
