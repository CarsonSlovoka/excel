package app

import (
	"github.com/CarsonSlovoka/excel/app/install"
	"log"
)

type about struct {
	ProgName string
	Version  string
	Author   string
}

var (
	About *about

	LoggerSys   *log.Logger
	LoggerError *log.Logger
)

func init() {
	install.InitConfigDir(ProgName)
}
