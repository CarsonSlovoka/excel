package urls

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Setting struct {
	EnableBootstrap bool
	EnableFontawesome bool
	EnableJquery    bool
}

var (
	SiteSetting Setting
)

func init() {
	SiteSetting = Setting{
		EnableBootstrap: true, EnableJquery: true, EnableFontawesome: true,
	}
}

type htmlTemplate struct {
	*template.Template
	context map[string]interface{}
}

func (t *htmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := t.Execute(w, t.context); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("BadRequest\n"))
		return
	}
	return
}

func NewTemplate(targetName string, fs fs.FS, patterns ...string) *htmlTemplate {
	ht, err := template.New(targetName).ParseFS(fs, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	return &htmlTemplate{ht, make(map[string]interface{})}
}
