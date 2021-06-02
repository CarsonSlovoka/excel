package urls

import (
    "embed"
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/CarsonSlovoka/excel/app"
    "github.com/CarsonSlovoka/excel/app/server"
    i18nPlugin "github.com/CarsonSlovoka/excel/pkg/i18n"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "log"
    "net/http"
    "path"
    "path/filepath"
    "regexp"
    "strings"
)

//go:embed i18n
var i18nDirFS embed.FS

type i18nObject struct {
    *i18n.Bundle
    messageFileMap map[string]*i18n.MessageFile
}

var i18nObj *i18nObject

func newI18nObj() *i18nObject {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    i18nDir := "i18n"

    dirEntrySlice, err := i18nDirFS.ReadDir(i18nDir)
    if err != nil {
        log.Fatal(err)
    }

    messageFileMap := make(map[string]*i18n.MessageFile, 0)
    for _, dirEntry := range dirEntrySlice {
        if dirEntry.IsDir() {
            continue
        }
        filename := dirEntry.Name()
        langFilePath := path.Join(i18nDir, filename)

        bytesLang, err := i18nDirFS.ReadFile(langFilePath)
        if err != nil {
            log.Fatal(err)
        }
        messageFile, err := bundle.ParseMessageFileBytes(bytesLang, langFilePath)
        messageFileMap[strings.TrimSuffix(filename, filepath.Ext(filename))] = messageFile
    }
    return &i18nObject{bundle, messageFileMap}
}

// init i18nObj
func init() {
    if i18nObj != nil {
        return
    }
    i18nObj = newI18nObj()
}

func initI18nJS() {
    langTmpl := &i18nPlugin.LangTmpl{Bundle: i18nObj.Bundle}
    expr := `var i18n = {
{{range .MessageSet}} {{.}}: "{{i18n . ""}}",
{{end}}
}`

    i18nRouter := server.Mux.PathPrefix("/i18n/").Subrouter()

    regex := regexp.MustCompile("/i18n/(?P<lang>[a-z]{2}|[a-z]{2}-[a-z]{2})/") // en, zh-tw
    langIndex := regex.SubexpIndex("lang")

    for targetLang, _ := range i18nObj.messageFileMap {
        i18nRouter.HandleFunc(fmt.Sprintf("/%s/", targetLang),
            func(writer http.ResponseWriter, request *http.Request) {
                writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")

                matchSlice := regex.FindStringSubmatch(request.URL.Path)
                if matchSlice == nil {
                    return
                }
                curLang := matchSlice[langIndex]
                messageFile := i18nObj.messageFileMap[curLang]
                var messageIDSet []i18nPlugin.MessageID
                for _, message := range messageFile.Messages {
                    messageIDSet = append(messageIDSet, i18nPlugin.MessageID(message.ID))
                }

                langTmpl.MustCompile(curLang, expr, map[string]interface{}{
                    "Version": app.Version, // i18n/en.toml
                    "Author":  app.Author,
                })

                // io.MultiWriter(writer, os.Stdout)
                langTmpl.MustRender(writer, i18nPlugin.Context{
                    "MessageSet": messageIDSet,
                })
            }).
            Methods("GET")
    }
}
