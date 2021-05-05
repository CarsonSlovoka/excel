package app

import (
    "github.com/CarsonSlovoka/excel/api/utils"
    "github.com/CarsonSlovoka/excel/app/install"
    "log"
    "os"
    "path/filepath"
    "time"
)

func init() {
    logSysFilePath := filepath.Join(install.Path.ConfigDir, "sys.log")
    logErrFilePath := filepath.Join(install.Path.ConfigDir, "error.log")
    for _, fPath := range []string{logErrFilePath, logSysFilePath} {

        fileInfo, err := os.Stat(fPath)
        if os.IsNotExist(err) {
            continue
        }
        // if time.Now().After(fileInfo.ModTime().AddDate(0, 0, 7)) {
        if time.Now().After(fileInfo.ModTime().Add(10*time.Second)) {
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
