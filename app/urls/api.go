package urls

import (
    "encoding/json"
    "errors"
    "github.com/CarsonSlovoka/excel/app/server"
    http2 "github.com/CarsonSlovoka/excel/pkg/net/http"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "net/http"
    "os"
    "path/filepath"
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
    funcName := strings.TrimPrefix(urlPath, "/api/")
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

        if val, err := getValue(mapVal, inputPara); err == nil {
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
    case "path/filepath/Glob":
        para, err := getFirstValue(mapVal, "para")
        if err != nil {
            http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
        }
        var dirInfo FileInfo
        if err := json.Unmarshal([]byte(para), &dirInfo); err != nil {
            http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
            return
        }
        ext := req.FormValue("ext")
        if ext == "" {
            ext = "woff,woff2"
        }
        // response, _ := http.Get(fmt.Sprintf("http://127.0.0.1:%s%s", app.Port, filePath)) // 返回整個HTML的頁面資料太大，不適合 /user/static/a8308a402c6f40f65ca720ef7c54ecf0/
        // bodyBytes, err := ioutil.ReadAll(response.Body)
        var extSlice []interface{}
        for _, v := range strings.Split(ext, ",") {
            extSlice = append(extSlice, v)
        }
        var fontsSlice []string
        err = filepath.Walk(dirInfo.Path + "/fonts", func(path string, info os.FileInfo, err error) error {
            if info.IsDir() {
                return nil
            }
            if utils.In(filepath.Ext(path)[1:], extSlice...) { // ttf, woff, woff2 ...
                fontsSlice = append(fontsSlice, filepath.Base(path))
            }
            return nil
        })
        if err != nil {
            http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        byteData, _ := json.Marshal(fontsSlice)
        _, _ = w.Write(byteData)

    default:
        http2.ShowErrorRequest(w, http.StatusBadRequest, "wrong function name!")
    }
}

func getValue(valMap map[string][]string, targetName string) ([]string, error) {
    if val, exists := valMap[targetName]; exists {
        if len(val) != 0 {
            return val, nil
        }
    }
    return nil, errors.New("value not found error")
}

func getFirstValue(valMap map[string][]string, targetName string) (string, error) {
    result, err := getValue(valMap, targetName)
    if err != nil {
        return "", err
    }
    return result[0], nil
}

func init() {
    apiRouter := server.Mux.PathPrefix("/api/").Subrouter()
    apiRouter.HandleFunc("/os/{funName}", funcHandler).Methods("GET", "POST")
    // apiRouter.HandleFunc("/{package:(?:os|path)}/{funName}", funcHandler).Methods("GET", "POST")

    apiRouter.HandleFunc("/path/filepath/{funName}", funcHandler).Methods("GET", "POST")
}
