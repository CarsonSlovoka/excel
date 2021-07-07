package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

type Setting struct {
    EnableBootstrap   bool
    EnableTachyons    bool
    EnableFontawesome bool
    EnableJquery      bool
    ShowNavbar        bool
    ShowFooter        bool
}

var (
    SiteSetting Setting
)

func init() {
    SiteSetting = Setting{
        EnableBootstrap: true, EnableTachyons: true,
        EnableJquery: true, EnableFontawesome: true,
        ShowNavbar: true, ShowFooter: true,
    }
}

type Context = map[string]interface{} // alias

var (
    BaseContext Context
)

func init() {
    BaseContext = map[string]interface{}{
        "Site": SiteSetting,
        // "TabIcon": "/static/app.ico", // use `/favicon.ico` to instead of it.
        "AppName": app.ProgName,
        "Author":  app.Author,
        "Version": app.Version,
        // Lang: // determine at getting request
    }
}

func getLangContext(r *http.Request) Context {
    ctx := Context{"Lang": "en-us"} // default

    for _, cookieName := range []string{server.CookieNameMap.Config} {
        queryMap, err := server.SafeCookie.GetSecureCookieValue(r, cookieName)
        if err != nil {
            return ctx
        }
        for key, val := range queryMap {
            ctx[key] = val
        }
    }

    return ctx
}
