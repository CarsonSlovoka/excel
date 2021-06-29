package urls

import (
    "encoding/json"
    "fmt"
    "github.com/CarsonSlovoka/excel/app/server"
    i18nPlugin "github.com/CarsonSlovoka/excel/pkg/i18n"
    http2 "github.com/CarsonSlovoka/excel/pkg/net/http"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "net/http"
    "strings"
    "time"
)

func init() {
    /*
       set lang from URL, for example:
        - /config/?lang=en
        - /config/?lang=zh-tw
    */
    server.Mux.HandleFunc("/config/",
        func(w http.ResponseWriter, r *http.Request) {
            var outputMsg []string
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            queryMap, err := server.SafeCookie.GetSecureCookieValue(r, server.CookieNameMap.Config)
            getLang := func() (lang string) {
                val, exists := queryMap["Lang"]
                var previousLang string
                if !exists {
                    // previousLang = strings.Split(r.Header.Get("Accept-Language"), ",")[0]
                    previousLang = "en-us"
                } else {
                    previousLang = val.(string)
                }

                lang = r.FormValue("lang")
                if lang == "" {
                    lang = previousLang
                }

                var legalLangSet []string
                for curLang, _ := range i18nObj.messageFileMap {
                    legalLangSet = append(legalLangSet, curLang)
                }
                if !(utils.In(lang, legalLangSet)) {
                    outputMsg = append(outputMsg,
                        fmt.Sprintf(`{{i18n "Err.InputLangNotExist" (dict "Lang" "%s") | safeHTML}}`, lang),
                    )
                    lang = previousLang // default language
                }
                return lang
            }
            lang := getLang()
            updateDataMap := map[string]interface{}{
                "Lang": lang,
            }
            expires := time.Now().AddDate(0, 1, 0)
            if err := server.SafeCookie.UpdateSecureCookie(w, r, server.CookieNameMap.Config, updateDataMap, "", &expires); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                _, _ = w.Write([]byte(err.Error() + "Please try again!"))
                return
            }

            outputData := struct {
                Lang string
                Msg  []string `json:",omitempty"`
            }{
                updateDataMap["Lang"].(string),
                outputMsg,
            }

            beautifulJsonByte, err := json.MarshalIndent(outputData, "", "  ")
            if err != nil {
                http2.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
                return
            }

            // _, _ = w.Write(beautifulJsonByte)
            langTmpl := &i18nPlugin.LangTmpl{Bundle: i18nObj.Bundle}
            expr := strings.ReplaceAll(string(beautifulJsonByte), "\\", "")
            langTmpl.MustCompile(lang, expr, Context{})
            langTmpl.MustRender(w, nil)
        }).Methods("GET")
}
