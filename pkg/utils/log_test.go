package utils

import (
    "log"
    "os"
    "testing"
)

func TestInitFileLogger(t *testing.T) {
    filepath := ".test.temp.log"
    file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        panic(err)
    }
    defer func() {
        _ = file.Close()
        if err = os.Remove(filepath); err != nil {
            panic(err)
        }
    }()
    logger := InitFileLogger(file, log.Ldate|log.Ltime|log.Lshortfile, true)
    logger.Println("...")
}
