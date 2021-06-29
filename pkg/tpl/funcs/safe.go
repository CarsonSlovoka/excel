package funcs

import (
    "github.com/spf13/cast"
    "html/template"
)

func CSS(a interface{}) (template.CSS, error) {
    s, err := cast.ToStringE(a)
    return template.CSS(s), err
}

func HTML(a interface{}) (template.HTML, error) {
    s, err := cast.ToStringE(a)
    return template.HTML(s), err
}

func HTMLAttr(a interface{}) (template.HTMLAttr, error) {
    s, err := cast.ToStringE(a)
    return template.HTMLAttr(s), err
}

func JS(a interface{}) (template.JS, error) {
    s, err := cast.ToStringE(a)
    return template.JS(s), err
}

func JSStr(a interface{}) (template.JSStr, error) {
    s, err := cast.ToStringE(a)
    return template.JSStr(s), err
}

func URL(a interface{}) (template.URL, error) {
    s, err := cast.ToStringE(a)
    return template.URL(s), err
}
