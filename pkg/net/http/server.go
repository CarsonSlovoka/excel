package http

import (
    "encoding/json"
    "net/http"
    "reflect"
)

func ShowErrorRequest(writer http.ResponseWriter, statusCode int, errMsg interface{}) {
    /*
       Usage:

           - ShowErrorRequest(w, 400, err.Error())
           - ShowErrorRequest(w, http.StatusBadRequest, map[string]interface{}{"Message": "錯誤!", "Time": time.Now()})
    */
    writer.WriteHeader(statusCode)
    if errMsg == nil {
        return
    }

    value := reflect.ValueOf(errMsg)
    switch value.Kind() {
    case reflect.String:
        _, _ = writer.Write([]byte(value.String()))
    case reflect.Struct:
        fallthrough
    case reflect.Map:
        writer.Header().Set("Content-Type", "application/json; charset=utf-8")
        if beautifulJsonByte, err := json.MarshalIndent(errMsg, "", "  "); err == nil {
            _, _ = writer.Write(beautifulJsonByte)
        }
    }
}
