package ck

import (
    "encoding/json"
    "errors"
    "github.com/stretchr/testify/assert"
    "log"
    "net/http"
    "testing"
    "time"
)

func write2JSON(w http.ResponseWriter, obj interface{}) {
    w.Header().Set("Content-Type", "application/json")
    beautifulJsonByte, err := json.MarshalIndent(obj, "", "  ")
    if err != nil {
        panic(err)
    }
    _, _ = w.Write(beautifulJsonByte)
}

var safeCookie *SafeCookie

func init() {
    hashKey := []byte("AJCMZ!@KDOADFJ!#$^*UPQ)@O#MDOFQVK$^&*DJQG<BOEDK^)(O20oda0khd=s2-")
    blockKey := []byte("&*DJQG<BOEDK^)(O20oda0khd=s2-$^*")
    if len(hashKey) != 64 && len(blockKey) != 32 {
        panic("key length error")
    }

    safeCookie = New(hashKey, blockKey,
        300) // 5 minute
}

func getCookies(url string) []*http.Cookie {
    response, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    return response.Cookies()
}

func getCookieByName(url, cookieName string) *http.Cookie {
    cookies := getCookies(url)
    for _, cookie := range cookies {
        if cookie.Name == cookieName {
            return cookie
        }
    }
    return nil
}

func TestBasic(t *testing.T) {

    const CookieName = "MyCookie"

    setHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ck := NewCookie(CookieName)
        ck.Path = "/"
        ck.Expires = time.Now().AddDate(0, 1, 0)

        safeCookie.SetSecureCookie(ck, map[string]interface{}{
            "desc": "test cookie",
            "msg":  "hello world",
        }, w)

    })

    http.HandleFunc("/set-no-redirect/", func(w http.ResponseWriter, r *http.Request) {
        setHandlerFunc(w, r)
    })

    http.HandleFunc("/set/", func(w http.ResponseWriter, r *http.Request) {
        setHandlerFunc(w, r)
        http.Redirect(w, r, "/get/", http.StatusSeeOther)
    })

    http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
        queryMap, err := safeCookie.GetSecureCookieValue(r, CookieName)
        if err != nil {
            _, _ = w.Write([]byte(err.Error()))
            return
        }
        write2JSON(w, queryMap)
    })

    http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
        ClearCookie(CookieName, w)
        http.Redirect(w, r, "/get/", http.StatusSeeOther)
    })

    go func() {
        _ = http.ListenAndServe(":8080", nil)
    }()

    cookie := getCookieByName("http://localhost:8080/set-no-redirect/", CookieName)
    if cookie == nil {
        log.Fatal(errors.New("empty cookie value"))
    }

    request, err := http.NewRequest("GET", "http://localhost:8080/get/", nil)
    if err != nil {
        panic(err)
    }
    request.AddCookie(cookie)
    // cookie.Value // That was encoded, so we need to use the GetSecureCookieValue to restore it.
    queryMap, err := safeCookie.GetSecureCookieValue(request, CookieName)
    assert.Equal(t, "test cookie", queryMap["desc"])

    select {
    case <-time.After(1 * time.Second):
        return
    }
}
