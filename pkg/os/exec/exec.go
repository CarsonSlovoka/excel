package exec

import (
    "fmt"
    "log"
    "os/exec"
    "runtime"
)

func OpenBrowser(url string) {
    // 黄勇刚(hyg) https://gist.github.com/hyg/9c4afcd91fe24316cbf0
    var err error

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatal(err)
    }
}
