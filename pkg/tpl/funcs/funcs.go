package funcs

import (
    "errors"
    "reflect"
)

func Dict(values ...interface{}) (map[string]interface{}, error) {
    if len(values)%2 != 0 {
        return nil, errors.New("parameters must be even")
    }
    dict := make(map[string]interface{})
    var key, val interface{}
    for {
        key, val, values = values[0], values[1], values[2:]
        switch reflect.ValueOf(key).Kind() {
        case reflect.String:
            dict[key.(string)] = val
        default:
            return nil, errors.New(`type must equal to "string"`)
        }
        if len(values) == 0 {
            break
        }
    }
    return dict, nil
}
