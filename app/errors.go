package app

import "errors"

var ErrCode errorCode

type errorCode struct {
    SecureCookieDecodeError error
}

func init() {
    ErrCode = errorCode{
        errors.New("secureCookie decode error"),
    }
}
