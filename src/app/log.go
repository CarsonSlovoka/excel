package app

import (
    "github.com/CarsonSlovoka/excel/app/install"
    osPlugin "github.com/CarsonSlovoka/excel/pkg/os"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "log"
    "os"
    "path/filepath"
    "time"
)

func init() {
    initLog()
}

func initLog() {
    logSysFilePath := filepath.Join(install.Path.ConfigDir, "log/sys.log")
    logErrFilePath := filepath.Join(install.Path.ConfigDir, "log/error.log")

    for _, fPath := range []string{logErrFilePath, logSysFilePath} {

        osPlugin.MustMkdirAll(filepath.Dir(fPath))

        fileInfo, err := os.Stat(fPath)
        if os.IsNotExist(err) {
            continue
        }
        if time.Now().After(fileInfo.ModTime().AddDate(0, 0, 7)) {
            // Remove the log file if it's not used for a while. We treat this as a not important log file.
            err = os.Remove(fPath)
            if err != nil {
                panic(err)
            }
        }
    }

    f, err := os.OpenFile(logSysFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) // os.O_WRONLY|os.O_CREATE move the cursor to the beginning and then write the data; hence it will not remove the old data.
    if err != nil {
        panic(err)
    }
    LoggerSys = utils.InitFileLogger(f, log.Ldate|log.Ltime, true)

    f, err = os.OpenFile(logErrFilePath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        panic(err)
    }
    LoggerError = utils.InitFileLogger(f, log.Ldate|log.Ltime|log.Lshortfile, true)
}
