package ck

import (
    "github.com/gorilla/securecookie"
    "net/http"
)

type SafeCookie struct {
    cookie *securecookie.SecureCookie
}

type CookieValueMap map[string]interface{}

func New(hashKey []byte, blockKey []byte, maxAge int) *SafeCookie {
    /*
       recommend
           hashKey.length: 64
           blockKey.length: 32

       maxAge: Default is 86400 * 30. 1 month
    */
    secureCookie := securecookie.New(hashKey, blockKey)
    secureCookie.MaxAge(maxAge)
    return &SafeCookie{secureCookie}
}

func (sc *SafeCookie) SetSecureCookie(cookie *http.Cookie, mapValue CookieValueMap, responseWriter http.ResponseWriter) {

    encoded, err := sc.cookie.Encode(cookie.Name, mapValue)
    if err == nil {
        encodeCookie := &http.Cookie{
            Name:     cookie.Name,
            Value:    encoded,
            Expires:  cookie.Expires, // If unspecified, the cookie becomes a session cookie. A session finishes when the client shuts down
            HttpOnly: cookie.HttpOnly,
            SameSite: cookie.SameSite,
            Path:     cookie.Path,
            MaxAge:   cookie.MaxAge,
        }
        http.SetCookie(responseWriter, encodeCookie)
    }
}

func (sc *SafeCookie) GetSecureCookieValue(r *http.Request, cookieName string) (cookieValue CookieValueMap, err error) {
    var ck *http.Cookie
    ck, err = r.Cookie(cookieName)
    // cookieValue := make(CookieValueMap)
    if ck != nil {
        err = sc.cookie.Decode(ck.Name,
            ck.Value,
            &cookieValue)
        return cookieValue, err
    }
    return
}
