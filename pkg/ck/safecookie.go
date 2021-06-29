package ck

import (
    "errors"
    "github.com/gorilla/securecookie"
    "net/http"
    "time"
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

// Get previous data if the value is nil will new an empty data.
func (sc *SafeCookie) GetSecureCookieValue(r *http.Request, cookieName string) (cookieValue CookieValueMap, err error) {
    var ck *http.Cookie
    ck, err = r.Cookie(cookieName)
    // cookieValue := make(CookieValueMap)
    if ck != nil {
        if err = sc.cookie.Decode(ck.Name,
            ck.Value,
            &cookieValue); err != nil {
            return nil, errors.New("secureCookie decode error")
        }
        return cookieValue, nil
    }
    return make(CookieValueMap), err
}

func (sc *SafeCookie) UpdateSecureCookie(writer http.ResponseWriter, request *http.Request, cookieName string,
    updateMap map[string]interface{},
    path string, expires *time.Time,
) error {
    /*
       expires= time.Now().AddDate(0, 1, 0)
    */
    cookieValueMap, err := sc.GetSecureCookieValue(request, cookieName)

    if cookieValueMap == nil {
        // This is a rare occurrence and may be the result of an abnormal termination procedure.
        ClearCookie(cookieName, writer)
        return err
    }

    for key, val := range updateMap {
        cookieValueMap[key] = val
    }

    cookie := NewCookie(cookieName)
    if path == "" {
        cookie.Path = "/"
    }
    if expires != nil {
        cookie.Expires = *expires
    }

    sc.SetSecureCookie(cookie, cookieValueMap, writer)
    return nil

}
