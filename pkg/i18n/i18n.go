package i18n

import (
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "html/template"
    "io"
    "log"
)

type LangTmpl struct {
    Bundle *i18n.Bundle
    *i18n.Localizer
    *template.Template
}

type MessageID string
type Context map[string]interface{}

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

func (i18nTmpl *LangTmpl) GetI18nFuncMap(templateData map[string]interface{}) template.FuncMap {
    i18nFunc := func(messageID MessageID) string {
        // {{ i18n "whatsInThis" . }}
        return i18nTmpl.Localizer.MustLocalize(&i18n.LocalizeConfig{
            MessageID:    string(messageID),
            TemplateData: templateData, // other = "What's in this {{ .Type }}"
        })
    }
    return template.FuncMap{"i18n": i18nFunc, "T": i18nFunc}
}

func (i18nTmpl *LangTmpl) Render(wr io.Writer, ctx Context) error {
    if err := i18nTmpl.Template.Execute(wr, ctx); err != nil {
        return err
    }
    return nil
}

func (i18nTmpl *LangTmpl) MustRender(wr io.Writer, ctx Context) {
    if err := i18nTmpl.Render(wr, ctx); err != nil {
        log.Fatal(err)
    }
}
