package funcs

import "html/template"

func GetUtilsFuncMap() template.FuncMap {
    i18nFunc := func(messageID string, templateData interface{}) string { return messageID } // Just let "i18n" and T is legal.
    return template.FuncMap{
        "i18n": i18nFunc, "T": i18nFunc,
        "dict":  Dict,
        "Slice": Slice, // Let 1st char uppercase since "slice" was defined already.
        "split": Split,
    }
}
