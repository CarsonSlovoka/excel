package server

import "github.com/CarsonSlovoka/excel/pkg/ck"

var (
    SafeCookie    *ck.SafeCookie
    CookieNameMap *cookieNameMap
)

type cookieNameMap struct {
    Config string
}

// cookie
func init() {
    CookieNameMap = &cookieNameMap{}
    CookieNameMap.Config = "config"

    hashKey := []byte("AJCMZ!@KDOADFJ!#$^*UPQ)@O#MDOFQVK$^&*DJQG<BOEDK^)(O20oda0khd=s2-")
    blockKey := []byte("&*DJQG<BOEDK^)(O20oda0khd=s2-$^*")
    if len(hashKey) != 64 && len(blockKey) != 32 {
        panic("key length error")
    }
    SafeCookie = ck.New(hashKey, blockKey, 86400)
}
