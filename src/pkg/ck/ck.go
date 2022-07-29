package ck

import "net/http"

func NewCookie(cookieName string) *http.Cookie {
    cookie := &http.Cookie{
        Name: cookieName,
    }
    return cookie
}

func ClearCookie(name string, rw http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   name,
        Value:  "",
        Path:   "/",
        MaxAge: -1, // delete cookie now, equivalently 'Max-Age: 0'
    }
    http.SetCookie(rw, cookie)
}
