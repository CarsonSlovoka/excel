package session

import (
    "context"
    "fmt"
    "github.com/stretchr/testify/assert"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "testing"
    "time"
)

var (
    SessionManager *Manager

    Mux *http.ServeMux
)

const (
    CookieName = "session"
)

func init() {
    var err error
    providerName := "my-memory-provider"
    Register(providerName, GetMemoryProvider())
    SessionManager, err = NewManager(providerName, CookieName, 500) // 5 minute
    if err != nil {
        panic(err)
    }
    go SessionManager.GC() // start GC Processing
}

func init() {
    Mux = http.NewServeMux()
}

func listenAndServe(quit chan bool) {
    server := http.Server{Addr: ":8080", Handler: Mux}

    Mux.HandleFunc("/shutdown/", func(w http.ResponseWriter, r *http.Request) {
        if err := server.Shutdown(context.Background()); err != nil {
            panic(err)
        }
    })

    if err := server.ListenAndServe(); err != nil {
        log.Println(err.Error())
    }

    quit <- true
}

func sendRequest(method, url string, reader io.Reader, cookies []*http.Cookie) string {
    request, err := http.NewRequest(method, url, reader)
    // request.Header.Set("Content-Type", "application/json")
    if err != nil {
        panic(err)
    }

    cookieValue := ""
    for _, cookie := range cookies {
        if cookie.Name != CookieName {
            continue
        }
        cookieValue = cookie.Value
        request.AddCookie(&http.Cookie{Name: CookieName, Value: cookieValue})
        break
    }

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    // response.Status

    defer func() { _ = response.Body.Close() }()
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }

    return string(bodyBytes)
}

func getRequest(url string) (string, []*http.Cookie) {
    response, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer func() { _ = response.Body.Close() }()
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }

    return string(bodyBytes), response.Cookies()
}

func TestBasic(t *testing.T) {

    Mux.Handle("/start-no-redirect/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        store := SessionManager.SessionStart(w, r)
        store = store.(*MemoryStore)
        if err := store.Set("msg", "Hello World"); err != nil {
            panic(err)
        }
        // http.Redirect(w, r, "/get/", http.StatusSeeOther) // use the http.StatusSeeOther can't get the cookie
    }))

    Mux.Handle("/start/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        store := SessionManager.SessionStart(w, r)
        store = store.(*MemoryStore)
        if err := store.Set("msg", "Hello World"); err != nil {
            panic(err)
        }
        http.Redirect(w, r, "/get/", http.StatusSeeOther)
    }))

    Mux.Handle("/set/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session := SessionManager.SessionQuery(r)
        if session == nil {
            fmt.Println("redirect to start")
            http.Redirect(w, r, "/start/", http.StatusSeeOther)
            return
        }

        store := session.(*MemoryStore)
        if err := store.Set("version", "0.0.0"); err != nil {
            panic(err)
        }
        http.Redirect(w, r, "/get/", http.StatusSeeOther)
        return
    }))

    Mux.Handle("/get/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session := SessionManager.SessionQuery(r)
        if session != nil {
            store := session.(*MemoryStore)
            outputString := ""
            for _, v := range store.GetMap() {
                outputString += v.(string) + " | "
            }
            _, _ = w.Write([]byte(outputString))
            return
        }
        _, _ = w.Write([]byte("empty"))
    }))

    Mux.Handle("/clear/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        SessionManager.SessionDestroy(w, r)
        http.Redirect(w, r, "/get/", http.StatusSeeOther)
    }))

    quit := make(chan bool)
    go listenAndServe(quit)

    _, cookies := getRequest("http://localhost:8080/start-no-redirect/")
    responseBody := sendRequest("GET", "http://localhost:8080/get/", nil, cookies)

    assert.Equal(t, "Hello World | ", responseBody)

    responseBody = sendRequest("GET", "http://localhost:8080/set/", nil, cookies)
    assert.Contains(t, responseBody, "Hello World")
    assert.Contains(t, responseBody, "0.0.0")

    responseBody = sendRequest("GET", "http://localhost:8080/clear/", nil, nil)
    assert.Equal(t, "empty", responseBody)

    select {
    case <-quit: // you can visit `http://localhost:8080/shutdown/` or wait 5 seconds to auto-finished.
        return
    case <-time.After(5 * time.Second):
        return
    }

}
