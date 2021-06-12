package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "github.com/CarsonSlovoka/excel/pkg/i18n"
    i18nPlugin "github.com/CarsonSlovoka/excel/pkg/i18n"
    httpPlugin "github.com/CarsonSlovoka/excel/pkg/net/http"
    "github.com/CarsonSlovoka/excel/pkg/tpl/funcs"
    "html/template"
    "io/fs"
    "log"
    "net/http"
)

type htmlTemplate struct {
    *i18n.LangTmpl // *template.Template
    contextSet     []Context
}

func (t *htmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := Context{}
    for _, curCtx := range t.contextSet {
        for k, v := range curCtx {
            ctx[k] = v
        }
    }

    noLangHandler := func() {
        if err := t.Execute(w, ctx); err != nil {
            log.Println(err)
            httpPlugin.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
            return
        }
    }

    queryMap, _ := server.SafeCookie.GetSecureCookieValue(r, server.CookieNameMap.Config)
    if queryMap == nil {
        app.LoggerError.Printf("Missing language: %s. Set it on URL: /config/?lang=", app.ErrCode.SecureCookieDecodeError.Error())
        noLangHandler()
        return
    }
    lang, exists := queryMap["Lang"]

    if !exists {
        noLangHandler()
        return
    }

    t.MustCompile(lang.(string), "", ctx)
    if err := t.Render(w, ctx); err != nil {
        log.Println(err)
        httpPlugin.ShowErrorRequest(w, http.StatusBadRequest, "BadRequest\n")
        return
    }
}

func NewTemplate(targetName string, fs fs.FS, patterns ...string) *htmlTemplate {
    if i18nObj == nil { // Because we can't make sure the init of I18n are done. If not it will be nil.
        i18nObj = newI18nObj()
    }
    tmplFuncs := funcs.GetUtilsFuncMap() // Adding extra FuncMap, by default LangTmpl, will also add it, but if you don't provide Language, that will render it with a standard HTML template; that is why are we need it.
    ht, err := template.New(targetName).Funcs(tmplFuncs).ParseFS(fs, patterns...)
    if err != nil {
        log.Fatal(err)
    }
    langTmpl := &i18nPlugin.LangTmpl{Bundle: i18nObj.Bundle, Template: ht}
    return &htmlTemplate{langTmpl, nil}
}

func ParseFiles(targetName string, filepath ...string) *htmlTemplate {
    if i18nObj == nil { // Because we can't make sure the init of I18n are done. If not it will be nil.
        i18nObj = newI18nObj()
    }

    langTmpl := &i18nPlugin.LangTmpl{Bundle: i18nObj.Bundle,
        Template: template.Must(template.New(targetName).Funcs(funcs.GetUtilsFuncMap()).ParseFiles(filepath...))}

    return &htmlTemplate{langTmpl, nil}
}
