package urls

import (
    "github.com/CarsonSlovoka/excel/app"
)

type Setting struct {
    EnableBootstrap bool
    EnableFontawesome bool
    EnableJquery    bool
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
