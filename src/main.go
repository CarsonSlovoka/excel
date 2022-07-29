package main

import (
	"fmt"
	"github.com/CarsonSlovoka/excel/app"
	"github.com/CarsonSlovoka/excel/app/server"
	"github.com/CarsonSlovoka/excel/app/urls"
	"github.com/CarsonSlovoka/excel/pkg/os/exec"
	"log"
	"sync"
	"time"
)

func main() {
	urls.InitURLs()
	urls.ShowAllURL()
	quit := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
		wg.Done()
		quit <- true
	}()
	time.Sleep(50 * time.Millisecond) // wait for listener ready
	go exec.OpenBrowser(fmt.Sprintf("http://127.0.0.1:%s", app.Port))
	wg.Wait()
	log.Println("Close App.")
}
