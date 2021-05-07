package utils

import (
    "errors"
    "github.com/stretchr/testify/assert" //
    "testing"
    "text/template"
)

func TestFString(t *testing.T) {
    // https://golang.org/pkg/text/template/#hdr-Actions
    // Example 1
    fs := &FString{}
    fs.MustCompile(`Name: {{.Name}} Msg: {{.Msg}}`, nil)
    fs.MustRender(map[string]interface{}{
        "Name": "Carson",
        "Msg":  123,
    })
    assert.Equal(t, "Name: Carson Msg: 123", fs.Data)
    fs.Clear()

    // Example 2 (with FuncMap)
    funcMap := template.FuncMap{
        "largest": func(slice []float32) float32 {
            if len(slice) == 0 {
                panic(errors.New("empty slice"))
            }
            max := slice[0]
            for _, val := range slice[1:] {
                if val > max {
                    max = val
                }
            }
            return max
        },
        "sayHello": func() string {
            return "Hello"
        },
    }
    fs.MustCompile("{{- if gt .Age 80 -}} Old {{else}} Young {{- end -}}"+ // "-" is for remove empty space
        "{{ sayHello }} {{largest .Numbers}}",
        funcMap)
    fs.MustRender(Context{
        "Age":     90,
        "Numbers": []float32{3, 9, 13.9, 2.1, 7},
    })
    assert.Equal(t, "Old Hello 13.9", fs.Data)
}
