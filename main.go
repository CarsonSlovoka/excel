package main

import (
	"fmt"
	"github.com/CarsonSlovoka/excel/app"
	"github.com/CarsonSlovoka/excel/app/server"
	"github.com/CarsonSlovoka/excel/app/urls"
	"github.com/CarsonSlovoka/excel/pkg/os/exec"
	"log"
)

func main() {
	urls.InitURLs()
	urls.ShowAllURL()
	quit := make(chan bool)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
		quit <- true
	}()

	go exec.OpenBrowser(fmt.Sprintf("http://localhost:%s", app.Port))

	for {
		select {
		case <-quit:
			log.Println("Close App.")
			return
		}
	}
}
