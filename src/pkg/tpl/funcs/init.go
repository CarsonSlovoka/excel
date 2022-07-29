package funcs

import "html/template"

func GetUtilsFuncMap() template.FuncMap {
    i18nFunc := func(messageID string, templateData interface{}) string { return messageID } // Just let "i18n" and T is legal. You can override it later.
    return template.FuncMap{
        "i18n": i18nFunc, "T": i18nFunc,
        "dict":         Dict,
        "Slice":        Slice, // Let 1st char uppercase since "slice" was defined already.
        "split":        Split,
        "replace":      Replace,
        "safeHTML":     HTML,
        "safeHTMLAttr": HTMLAttr,
        "safeCSS":      CSS,
        "safeJS":       JS,
        "safeJSStr":    JSStr,
        "safeURL":      URL,

        // 👇 Math
        "add":          Add,
        "sub":          Sub,
        "mul":          Mul,
        "div":          Div,
        "ceil":         Ceil,
        "floor":        Floor,
        "log":          Log,
        "sqrt":         Sqrt,
        "mod":          Mod,
        "modBool":      ModBool,
        "pow":          Pow,
        "round":        Round,
    }
}
