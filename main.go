package main

import (
    "github.com/CarsonSlovoka/excel/app/server"
    "github.com/CarsonSlovoka/excel/app/urls"
    "log"
)

func main() {
    urls.InitURLs()
    urls.ShowAllURL()
    log.Fatal(server.ListenAndServe())
}
