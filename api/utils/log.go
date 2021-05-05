package utils

import (
    "io"
    "log"
    "os"
)

func InitFileLogger(file *os.File, flag int, console bool) *log.Logger {
    writers := []io.Writer{
        file,
    }
    if console {
        writers = append(writers, os.Stdout)
    }
    // io.MultiWriter(file, os.Stdout)
    multipleWriter := io.MultiWriter(writers...) // same as above
    prefixString := ""
    return log.New(multipleWriter, prefixString, flag)
}
