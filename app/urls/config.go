package urls

import (
    "encoding/json"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
    "strings"
    "time"
)

func init() {
    /*
       set lang from URL, for example:
        - /config/?lang=en
        = /config/?lang=zh-tw
    */
    server.Mux.HandleFunc("/config/",
        func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json; charset=utf-8")

            queryMap, err := server.SafeCookie.GetSecureCookieValue(r, server.CookieNameMap.Config)

            lang := r.FormValue("lang")
            if lang == "" {
                val, exists := queryMap["Lang"]
                if !exists {
                    lang = strings.Split(r.Header.Get("Accept-Language"), ",")[0]
                } else {
                    lang = val.(string)
                }
            }

            outputData := map[string]interface{}{
                "Lang": lang,
            }
            expires := time.Now().AddDate(0, 1, 0)
            if err := server.SafeCookie.UpdateSecureCookie(w, r, server.CookieNameMap.Config, outputData, "", &expires); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                _, _ = w.Write([]byte(err.Error()))
                return
            }
            beautifulJsonByte, err := json.MarshalIndent(outputData, "", "  ")
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                _, _ = w.Write([]byte(err.Error()))
                return
            }
            _, _ = w.Write(beautifulJsonByte)
        }).Methods("GET")
}
