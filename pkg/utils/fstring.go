package utils

import (
    "text/template"
)

type Context map[string]interface{}

type FString struct {
    Data     string
    template *template.Template
}

func (fs *FString) MustCompile(expr string, funcMap template.FuncMap) {
    fs.template = template.Must(template.New("f-string").
        Funcs(funcMap).
        Parse(expr))
}

func (fs *FString) Write(b []byte) (n int, err error) {
    fs.Data += string(b)
    return len(b), nil
}

func (fs *FString) Render(context map[string]interface{}) error {
    if err := fs.template.Execute(fs, context); err != nil {
        return err
    }
    return nil
}

func (fs *FString) MustRender(context Context) {
    if err := fs.Render(context); err != nil {
        panic(err)
    }
}

func (fs *FString) Clear() string {
    // return the data and clear it
    out := fs.Data
    fs.Data = ""
    return out
}
