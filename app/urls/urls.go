package urls

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type htmlTemplate struct {
	*template.Template
	contextSet []Context
}

func (t *htmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := Context{}
    for _, curCtx := range t.contextSet {
        for k, v := range curCtx {
            ctx[k] = v
        }
    }
	if err := t.Execute(w, ctx); err != nil {
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
	return &htmlTemplate{ht, nil}
}
