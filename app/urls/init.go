package urls

import (
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    "github.com/gorilla/mux"
)

func InitURLs() {
    initStaticFS() // set static dir
    serveSingleFile("/favicon.ico", "static/app.ico")
    initSystemURL()
    initHomeURL()
    initFileURL()
}

func ShowAllURL() {
    err := server.Mux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        var urlPath string
        var err error
        urlPath, err = route.GetPathTemplate()
        if err != nil {
            return err
        }
        app.LoggerSys.Printf("%v\n", urlPath)
        return nil
    })
    if err != nil {
        panic(err)
    }
}
