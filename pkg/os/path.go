package os

import "os"

func MustMkdirAll(filepath string) {
    //  MustMkdirAll(filepath.Join(os.Getenv("programdata"), "xxx/dir"))
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        err := os.MkdirAll(filepath, os.ModePerm)
        if err != nil {
            panic(err)
        }
    }
}
