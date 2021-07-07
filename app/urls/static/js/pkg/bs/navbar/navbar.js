/**
 * @param {string} text
 * @param {Object} options
 */
export function Item(text, options = {}) {
  this.text = text

  // 👇 options
  this.icon = options.icon ?? "" // fas a-times-circle fa-2x
  this.url = options.url ?? ""
  this.className = options.className ?? "" // active, disable
  this.callback = options.callback ?? (() => {
  })
}

export class Navbar {

  /**
   * @param {HTMLElement} targetNode
   * @param {Object} options
   */
  constructor(targetNode, options = {}) {
    this.navNode = targetNode
    this.options = options
    this.options = {
      style: {
        hoverClass: options.style?.hoverClass ?? "" // hover-light-blue
      },
      title: {
        text: options.title?.text ?? "",
        faIcon: { // fontawesome icon
          name: options.title?.faIcon?.name ?? undefined, // fa-file-excel
          size: options.title?.faIcon?.size ?? 2, // fa-2x
          color: options.title?.faIcon?.color ?? "initial"
        }
      },
    }

    this.initCSSStyle()
    this.initNavbar()
  }

  initCSSStyle() {
    const range = document.createRange()
    const frag = range.createContextualFragment(`
<style>
body {min-height: 100rem; padding-top: 4.5rem;}
nav {transition: top 0.5s;}
</style>
`)
    document.querySelector("head").append(frag)
  }

  initNavbar() {
    const range = document.createRange()
    const fragNav = range.createContextualFragment(
      `<nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
<div class="container-fluid col-md-10 offset-md-1"></div>
</nav>`
    )
    const container = fragNav.querySelector(`div[class^="container-fluid"]`)

    const initTitle = () => {
      const range = document.createRange()
      const titleObj = this.options.title
      const faIcon = titleObj.faIcon
      const titleIconText = faIcon.name === undefined ? "" : `<i class="far ${faIcon.name} fa-${faIcon.size}x me-2" style="color:${faIcon.color}"></i>`
      const fragTitle = range.createContextualFragment(titleIconText + `
<a class="navbar-brand" href="#"><b>${titleObj.text}</b></a>
<button class="navbar-toggler" type="button" data-bs-toggle="collapse"
data-bs-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
<span class="navbar-toggler-icon"></span>
</button>
`
      )

      container.append(fragTitle)
    }

    const initBody = () => {
      const range = document.createRange()
      const fragBody = range.createContextualFragment(`
<div class="collapse navbar-collapse" id="navbarCollapse">
    <ul class="navbar-nav me-auto mb-2 mb-md-0"></ul>
    <ul class="navbar-nav mb-2 mb-md-0"></ul>
    <form class="d-flex"></form>
</div>`
      )
      container.append(fragBody)
    }

    initTitle()
    initBody()

    const oldNavClass = this.navNode.className
    this.navNode.parentNode.replaceChild(fragNav, this.navNode) // this.navNode.replaceWith(fragNav)
    this.navNode = document.querySelector(`nav[class^="navbar"]`)
    this.navNode.className += " " + oldNavClass
  }


  /**
   * @param {Item} item
   * @return {DocumentFragment}
   */
  newItem(item) {
    const range = document.createRange()
    const href = item.url === "" ? "" : `href="${item.url}"`
    const className = item.className === "" ? `tc near-white ${this.options.style.hoverClass}` : item.className
    const faIcon = item.icon === "" ? "" : `<i class="${item.icon}"></i><br>`
    const frag = range.createContextualFragment(
      `<li class="nav-item"><a class="nav-link" ${href}><p class="${className}">${faIcon}${item.text}</p></a></li>`
    )

    frag.querySelector("p").onclick = item.callback
    return frag
  }

  /**
   * @param {string} btnName
   * @param {string} btnIcon
   * @param {Item[]} items
   * @param {string} className
   * @return {DocumentFragment}
   */
  newDropdownItem(btnName, btnIcon, items, className) {
    const range = document.createRange()
    btnIcon = btnIcon === "" ? "" : `<i class="${btnIcon}"></i><br>`
    const fragDropdown = range.createContextualFragment(`
<li class="nav-item">
<div class="tc dropdown">
  <button class="btn btn-dark dropdown-toggle${" " + (className ?? this.options.style.hoverClass)}" type="button" data-bs-toggle="dropdown" aria-expanded="false">${btnIcon}${btnName}</button>
  <ul class="dropdown-menu"></ul>
</div>
</li>
`
    )

    const menu = fragDropdown.querySelector(`ul[class^="dropdown-menu"]`)
    items.forEach(item => {
      const range = document.createRange()
      const href = item.url === "" ? "" : `href="${item.url}"`
      const faIcon = item.icon === "" ? "" : `<i class="far ${item.icon}"></i><br>`
      const fragItem = range.createContextualFragment(
        `<li><a class="dropdown-item ${item.className}" ${href}>${faIcon} ${item.text}</a></li>`
      )
      fragItem.querySelector("a").onclick = item.callback
      menu.append(fragItem)
    })
    return fragDropdown
  }

  /**
   * @param {Item|Item[]} item
   * @param {"left"|"right"} align
   * @param {Object} options
   */
  AppendItem(item, align = "left", options = {}) {
    const node = (align === "left" ?
        this.navNode.querySelector(`ul[class^="navbar-nav me-auto"]`) :
        this.navNode.querySelector(`ul[class="navbar-nav mb-2 mb-md-0"]`)
    )
    if (Array.isArray(item)) {
      const btnName = options.dropdown?.name ?? ""
      const btnIcon = options.dropdown?.icon ?? ""
      node.append(this.newDropdownItem(btnName, btnIcon, item, options.dropdown?.className))
      return
    }
    node.append(this.newItem(item))
  }

  /**
   * @param {string} htmlString
   */
  AddItem2Form(htmlString) {
    const node = this.navNode.querySelector(`form[class="d-flex"]`)
    const range = document.createRange()
    const fragItem = range.createContextualFragment(`${htmlString}`)
    node.append(fragItem)
  }

  EnableSmartHidden() {
    const controller = new AbortController()
    this.EnableSmartHidden.prevScrollPos = 0
    this.EnableSmartHidden.controller = controller
    document.addEventListener("scroll", () => {
      this.navNode.style.top = this.EnableSmartHidden.prevScrollPos > window.pageYOffset ? "0" : "-200px"
      this.EnableSmartHidden.prevScrollPos = window.pageYOffset
    }, {
      signal: controller.signal,
      once: false
    })
  }

  DisableSmartHidden() {
    this.EnableSmartHidden.controller.abort()
  }
}
