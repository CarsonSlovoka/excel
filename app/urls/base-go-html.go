package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "net/http"
)

type Setting struct {
    EnableBootstrap   bool
    EnableFontawesome bool
    EnableJquery      bool
}

var (
    SiteSetting Setting
)

func init() {
    SiteSetting = Setting{
        EnableBootstrap: true, EnableJquery: true, EnableFontawesome: true,
    }
}

type Context map[string]interface{}

var (
    BaseContext Context
)

func init() {
    BaseContext = map[string]interface{}{
        "Site": SiteSetting,
        // "TabIcon": "/static/app.ico", // use `/favicon.ico` to instead of it.
        "Author":  app.Author,
        "Version": app.Version,
    }
}

func getLangContext(r *http.Request) Context {
    ctx := Context{"Lang": "en"} // default

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
