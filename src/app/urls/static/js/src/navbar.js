import * as navbar from "../pkg/bs/navbar/navbar.js"

(
  () => {
    const nav = new navbar.Navbar(document.querySelector("nav"),
      {
        style: {
          hoverClass: "hover-light-blue",
        },
        title: {
          text: "Green Viewer",
          faIcon: {
            name: "fa-file-excel",
            size: 2,
            color: "#adff2f"
          }
        },
      }
    )
    nav.AppendItem(new navbar.Item(i18n.About, {icon: "fas fa-info-circle fa-2x", url: "/about/"}))
    const changeLangFunc = async (lang) => {
      const response = await fetch(`/config/?lang=${lang}`, {})
      if (!response.ok) {
        return
      }
      document.location.reload()
    }
    const dropdownLang = [
      new navbar.Item(i18n.LangEN_US, {
        callback: async () => {
          await changeLangFunc("en-us")
        }
      }),
      new navbar.Item(i18n.LangZH_TW, {
        callback: async () => {
          await changeLangFunc("zh-tw")
        }
      }),
    ]
    nav.AppendItem(dropdownLang, "left", {
      dropdown: {
        name: i18n.LabelLang,
        icon: "fas fa-globe fa-2x",
        // className: "hover-light-red"
      }
    })
    nav.AppendItem(
      new navbar.Item(i18n.LabelCloseApp, {
        url: "/shutdown/",
        icon: "far fa-times-circle fa-2x",
        className: "tc near-white hover-light-red"
      }),
      "right")
    nav.EnableSmartHidden()
  }
)()
