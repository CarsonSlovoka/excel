package install

import (
    "log"
    "os"
    "path/filepath"
    "strings"
)

type Paths struct {
	ConfigDir    string
}

var (
	Path Paths
)


// A private directory that belongs to the program and is used only the program only.
// `.[AppName]` directory
func InitConfigDir(dirName string) {
	configDir, err := filepath.Abs("." + strings.Split(dirName, ".")[0])
	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	Path.ConfigDir = configDir
}
