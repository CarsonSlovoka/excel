package utils

import (
    "strings"
)

func RFindAndReplace(strData, subStr, newStr string) string {
    if lastIndex := strings.LastIndex(strData, subStr); lastIndex > -1 {
        strData = strData[:lastIndex] + newStr + strData[lastIndex+1:]
    }
    return strData
}
