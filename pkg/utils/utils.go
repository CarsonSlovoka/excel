package utils

import (
    "math"
    "reflect"
)

func In(e interface{}, set ...interface{}) bool {
    isNaN := false
    switch eValue := reflect.ValueOf(e); eValue.Kind() {
    case reflect.Slice:
        panic("type error. you can not set the slice of type to the element.")
    case reflect.Float64:
        isNaN = math.IsNaN(eValue.Float())
    }

    for _, item := range set {
        itemValue := reflect.ValueOf(item)
        if isNaN && itemValue.Kind() == reflect.Float64 && math.IsNaN(itemValue.Float()) {
            return true
        }
        if e == item {
            return true
        }
    }
    return false
}
