package i18n

import (
    "github.com/CarsonSlovoka/excel/pkg/tpl/funcs"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "html/template"
    "io"
    "log"
    "reflect"
)

type LangTmpl struct {
    Bundle *i18n.Bundle
    *i18n.Localizer
    *template.Template
}

type MessageID = string
type context map[string]interface{}

func (i18nTmpl *LangTmpl) mustLegalLang(lang string) {
    _, _, err := language.ParseAcceptLanguage(lang)
    if err != nil {
        log.Fatal(err)
    }
}

func (i18nTmpl *LangTmpl) MustCompile(
    lang string,
    expr string,
    templateData map[string]interface{},
) {
    i18nTmpl.mustLegalLang(lang)
    i18nTmpl.Localizer = i18n.NewLocalizer(i18nTmpl.Bundle, lang)

    if i18nTmpl.Template == nil {
        // used for js see: https://github.com/CarsonSlovoka/excel/blob/e6c3599/app/urls/static/js/i18n/en.js
        i18nTmpl.Template = template.Must(template.New("").
            Funcs(i18nTmpl.GetI18nFuncMap(templateData)).
            Parse(expr))
    } else {
        // used for .go.html. the content should determine on the outside (I means it without "expr")
        i18nTmpl.Template = i18nTmpl.Template.Funcs(i18nTmpl.GetI18nFuncMap(templateData)) // ``i18n`` or ``T`` are allow // i18nTmpl.text.execFuncs
    }
}

// FuncMap: i18n, and other useful functions
func (i18nTmpl *LangTmpl) GetI18nFuncMap(templateData map[string]interface{}) template.FuncMap {
    i18nFunc := func(messageID interface{}, options interface{}) (string, error) {
        // {{ i18n "whatsInThis" . }}

        if options != nil {
            switch options.(type) {
            case context:
            case interface{}:
                if reflect.ValueOf(options).Kind() == reflect.Map {
                    for key, val := range options.(map[string]interface{}) {
                        templateData[key] = val
                    }
                }
            }
        }

        return i18nTmpl.Localizer.MustLocalize(&i18n.LocalizeConfig{
            MessageID:    messageID.(string),
            TemplateData: templateData, // other = "What's in this {{ .Type }}"
        }), nil
    }

    resultMap := template.FuncMap{} // add other useful functions
    for k, v := range funcs.GetUtilsFuncMap() {
        resultMap[k] = v
    }
    resultMap["i18n"] = i18nFunc
    resultMap["T"] = i18nFunc
    return resultMap
}

func (i18nTmpl *LangTmpl) Render(wr io.Writer, ctx context) error {
    if err := i18nTmpl.Template.Execute(wr, ctx); err != nil {
        return err
    }
    return nil
}

func (i18nTmpl *LangTmpl) MustRender(wr io.Writer, ctx context) {
    if err := i18nTmpl.Render(wr, ctx); err != nil {
        log.Fatal(err)
    }
}
