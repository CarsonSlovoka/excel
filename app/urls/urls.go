package urls

import (
    "github.com/CarsonSlovoka/excel/app/server"
    "github.com/CarsonSlovoka/excel/pkg/i18n"
    i18nPlugin "github.com/CarsonSlovoka/excel/pkg/i18n"
    httpPlugin "github.com/CarsonSlovoka/excel/pkg/net/http"
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

    noLangHandler :=  func() {
        if err := t.Execute(w, ctx); err != nil {
            log.Println(err)
            httpPlugin.ShowErrorRequest(w, http.StatusBadRequest, err.Error())
            return
        }
    }

    queryMap, err := server.SafeCookie.GetSecureCookieValue(r, server.CookieNameMap.Config)
    if err != nil {
        noLangHandler()
        return
    }
    lang, exists := queryMap["Lang"]

    if !exists {
        noLangHandler()
        return
    }

    t.MustCompile(lang.(string), "", ctx)
    if err := t.Render(w, i18nPlugin.Context(ctx)); err != nil {
        log.Println(err)
        httpPlugin.ShowErrorRequest(w, http.StatusBadRequest, "BadRequest\n")
        return
    }
}

func NewTemplate(targetName string, fs fs.FS, patterns ...string) *htmlTemplate {
    if i18nObj == nil { // Because we can't make sure the init of I18n are done. If not it will be nil.
        i18nObj = newI18nObj()
    }
    tmplFuncs := func() template.FuncMap {
        i18nFunc := func(messageID string) string {return messageID} // Just let "i18n" and T is legal. Don't worry. The implementation for the function will change when doing Compile.
        return template.FuncMap{"i18n": i18nFunc, "T": i18nFunc}
    }
    ht, err := template.New(targetName).Funcs(tmplFuncs()).ParseFS(fs, patterns...)
    if err != nil {
        log.Fatal(err)
    }
    langTmpl := &i18nPlugin.LangTmpl{Bundle: i18nObj.Bundle, Template: ht}
    return &htmlTemplate{langTmpl, nil}
}
