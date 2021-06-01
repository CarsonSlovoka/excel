package urls

import (
    "embed"
    "github.com/CarsonSlovoka/excel/app/server"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "io/fs"
    "net/http"
    "path/filepath"
)

//go:embed static
var staticFS embed.FS

type StaticFS struct {
    embed.FS
}

func (sfs StaticFS) Open(path string) (fs.File, error) {
    if utils.In(filepath.Ext(path), []string{".rst", ".md"}) {
        return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
    }

    return sfs.FS.Open(path)
}

func initStaticFS() {
    // <link rel="stylesheet" href="/static/css/style.css" type="text/css">
    server.Mux.PathPrefix("/static/").Handler( // Remember to use the `PathPrefix` so that the sub-path still works.
        http.FileServer(http.FS(StaticFS{staticFS})))
}

type SingleFS struct {
    embed.FS
    filepath string
}

func serveSingleFile(pattern string, filepath string) {
    server.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        http.FileServer(http.FS(SingleFS{staticFS, filepath})).ServeHTTP(w, r)
    })
}

func (sfs SingleFS) Open(string) (fs.File, error) {
    return sfs.FS.Open(sfs.filepath)
}
