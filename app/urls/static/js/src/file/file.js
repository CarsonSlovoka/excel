import * as base64 from "../../pkg/encoding/base64.js"
import {ReadFile} from "../../pkg/io/ioutil.js"

const IFRAME_ID = "iframe-bootstrap-table"
const UNIQUE_ID = "__id__"
const BUTTONS_CLASS = "primary"

function removeExtraSpace(stringData) {
  stringData = stringData.replace(/,( *)/gm, ",")  // remove extra space
  stringData = stringData.replace(/^ *| *$/gm, "") // remove space on the beginning and end.
  return stringData
}

function save2Server() {
  const inputFile = document.getElementById('uploadFile')
  const inputValue = inputFile.value
  if (inputValue === "") {
    return
  }

  const filename = inputValue.substring(inputValue.lastIndexOf('\\') + 1).split(".")[0] // get basename without extension
  const filenameB64 = base64.EncodeToString(filename) // Ensure Chinese filename is working.

  const xhr = new XMLHttpRequest()
  xhr.open("POST", filenameB64, true)

  xhr.onload = (oEvent) => {
    if (xhr.status === 200) {
    } else {
      alert(`${xhr.responseText}`)
    }
  }
  const bsTable = document.getElementById(IFRAME_ID).contentWindow.table
  const dataArray = bsTable.bootstrapTable('getData', {
    useCurrentPage: false, // all page
    includeHiddenRows: true,
    unfiltered: true, // include all data (unfiltered).
    // formatted:
  })
  const jsonString = JSON.stringify(dataArray)
  const formData = new FormData() // https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects
  formData.set("uploadData", jsonString)
  xhr.send(formData)
}


/**
 * @param {string} url
 * @param {FileInfo} dirInfo
 */
function StaticInfo(url, dirInfo) {
  this.url = url
  this.dirInfo = dirInfo
}

function FileInfo(name, size, modTime, isDir, path) {
  this.Name = name
  this.Size = size
  this.ModTime = modTime
  this.IsDir = isDir
  this.Path = path
}

async function loadFonts(aliasName, url) {
  /*
  * https://developer.mozilla.org/en-US/docs/Web/API/FontFace/FontFace
  * https://stackoverflow.com/questions/11355147/font-face-changing-via-javascript
  * */
  const font = new FontFace(aliasName, `url(${url})`)
  await font.load()
  return font
}

function queryCellStyle(cellStyle, attr, defaultVal) {
  if (cellStyle === undefined) {
    return defaultVal === undefined ? undefined : defaultVal
  }
  const val = cellStyle.css[attr]
  return val === undefined ? defaultVal : val
}

class BSTable {

  /**
   * @param {Array} dataArray
   * @param {StaticInfo} staticInfo
   * @param {string} locale
   */
  constructor(dataArray, staticInfo, locale) {
    return (async () => {
      this.columns = [] // means: bootstrap-table columns
      this.dataArray = dataArray
      this.staticInfo = staticInfo
      this.fontList = await this.getFonts()
      const div = document.getElementById('div-csv-data')
      div.querySelectorAll('*').forEach(node => node.remove())
      this.table = undefined
      this.iframe = this.newIframe()
      this.locale = locale
      this.options = {
        filter: {
          enable: true,
          isNeedInit: true,
        }
      }
      this.initIframeEvent()
      div.append(this.iframe)
    })()
  }

  /**
   * @return {Array}
   */
  async getFonts() {
    if (this.staticInfo.url === "") {
      return []
    }
    const formData = new FormData()
    formData.set("para", JSON.stringify(this.staticInfo.dirInfo))
    const response = await fetch("/api/path/filepath/Glob?ext=ttf,woff,woff2", {
      method: "post",
      body: formData
    })
    if (!response.ok) {
      const errMsg = await response.text()
      throw Error(`${response.statusText} (${response.status}) | ${errMsg} `)
    }

    return await response.json()
  }

  newIframe() {
    const iframe = document.createElement('iframe')
    iframe.id = IFRAME_ID
    iframe.height = "748px"
    iframe.width = "100%"
    iframe.src = "/bs-table/"
    return iframe
  }

  getColumn(fieldName) {
    for (const col of this.columns) {
      if (col.field === fieldName) {
        return col
      }
    }
    return undefined
  }

  updateFilter() {
    if (!this.options.filter.enable || this.options.filter.isNeedInit) {
      throw Error("Filter Need Init or Enable")
    }
    this.table.bootstrapTable("filterBy", {
        "items": this.iframe.contentWindow.multipleSelect.multipleSelect('getSelects', "text"), // https://multiple-select.wenzhixin.net.cn/docs/en/methods#getselects
      }, {
        "filterAlgorithm": (rowObj, options = {"items": []}) => {
          const items = options.items
          if (items.length === 0) { // select all
            return true
          }
          const filterObj = {}
          for (const curColSelectString of items) {
            const { groups: { fieldName, valString } } = /(\[(?<fieldName>.*): (?<valString>.*),?)\]/gm.exec(`${curColSelectString}`)
            filterObj[fieldName] = valString.split(",").map(e=>e.trim())
          }

          for (const [fieldName, value] of Object.entries(rowObj)) {
            if (value === undefined) {
              continue
            }
            const okArray = filterObj[fieldName]
            if (okArray === undefined) {
              continue
            }
            if (okArray.includes(value)) {
              return true
            }
          }
          return false
        }
      }
    )
  }

  /**
   * @param {Object} targetCol is one of the elements of this.columns
   * @param {Object} cssObj
   * @example
   *  - updateBSTableColumnStyle("Name", {"font-family": myFont, "font-weight": 900, "background-color":"#da1235"})
   */
  updateBSTableColumnStyle(targetCol, cssObj = {}) {
    if (Object.keys(cssObj).length > 0) {
      const oldCellStyle = targetCol["cellStyle"]
      let newCSSObj = {css: {}}
      if (oldCellStyle !== undefined) {
        newCSSObj = oldCellStyle()
      }
      for (const [attr, value] of Object.entries(cssObj)) {
        newCSSObj.css[attr] = value
      }
      targetCol["cellStyle"] = () => {
        return newCSSObj // return {css: {}}
      }
    }

    const imgInfo = targetCol.imgInfo
    if (imgInfo.isImg) {
      targetCol.editable = undefined
      targetCol.formatter = (value, row, index, field) => {
        if (value.startsWith('http')) {
          return `<img src="${value}" alt="\???" style="width:${imgInfo.width};height:${imgInfo.height};" loading="lazy"/>`
        }
        return `<img src="${this.staticInfo.url}${value}" alt="\???" style="width:${imgInfo.width};height:${imgInfo.height};" loading="lazy"/>`
      }
    } else {
      targetCol.editable = {
        type: 'text',
        title: targetCol.field,
        emptytext: "",
        validate: function (v) {
          if (!v) return "You can't set the null value on this column"
        },
      }
      targetCol.formatter = (value, row, index, field) => {
        return value
      }
    }

    const hiddenColumns = this.table.bootstrapTable('getHiddenColumns')
    this.table.bootstrapTable('refreshOptions',
      {
        columns: this.columns,
      }
    )
    hiddenColumns.forEach((e) => {
      this.table.bootstrapTable('hideColumn', e.field)
    })
    if (this.options.filter.enable) { // The order is important. Do not put it in front of the "hideColumn"
      this.updateFilter()
    }
    this.table.bootstrapTable('refresh')
  }

  getCellStyle(fieldName) {
    for (let col of this.columns) {
      if (col.field === undefined || col.field !== fieldName) {
        continue
      }
      const cellStyleFunc = col["cellStyle"]
      if (cellStyleFunc === undefined) {
        return undefined
      }
      return cellStyleFunc()
    }
  }

  setConfigColumn(fieldName, titleName) {

    const iframeCtxWindow = this.iframe.contentWindow
    const iframeDoc = iframeCtxWindow.document

    iframeCtxWindow.showPopupConfig.onclick = (args) => { // add attribute for the function.
      const [fieldName, titleName] = args.split(",")
      const curColumn = this.getColumn(fieldName)
      const modal = iframeDoc.getElementsByClassName("popup-modal")[0]

      iframeDoc.getElementById("popup-modal-title").innerText = titleName

      const modalBody = iframeDoc.getElementsByClassName("popup-modal-body")[0]
      modalBody.querySelectorAll('*').forEach(node => node.remove())

      const divFont = document.createElement("div")
      const columnStyle = this.getCellStyle(fieldName)

      if ("Fonts Settings") {

        let oldFontFamily = queryCellStyle(columnStyle, "font-family", "????")
        let oldFontSize = queryCellStyle(columnStyle, "font-size", 8)
        divFont.className = "row"
        divFont.innerHTML = `
<select dir="rtl" class="pe-5 col-md-8 form-select" aria-label="select fonts" style="font-size:2em;">
<option selected>${oldFontFamily}</option>
</select>

<select dir="rtl" class="ms-2 pe-5 col-md-3 form-select" aria-label="select font-size">
<option selected>${oldFontSize}</option>
</select
`
        // SIZE
        const selectSize = divFont.querySelector(`select[aria-label="select font-size"]`)
        const N = 128, sizeArray = Array(N)
        for (let i = 0; i < N;) {
          sizeArray[i++] = i + 8
        }
        for (const curSize of sizeArray) {
          const nodeOption = document.createElement("option")
          nodeOption.value = curSize
          nodeOption.innerText = curSize
          selectSize.append(nodeOption)
        }
        selectSize.onchange = (e) => {
          const select = e.target
          const sizeValue = select.options[select.selectedIndex].value
          this.updateBSTableColumnStyle(curColumn, {"font-size": `${sizeValue}px`})
        }

        // FONT FAMILY
        const selectFonts = divFont.querySelector(`select[aria-label="select fonts"]`)
        for (const curFont of this.fontList) {
          const nodeOption = document.createElement("option")
          nodeOption.value = curFont
          nodeOption.innerText = curFont
          selectFonts.append(nodeOption)
        }
        selectFonts.onchange = async (e) => {
          const select = e.target
          const curFont = select.options[select.selectedIndex].value
          const fontAlias = curFont.replace(/\..*$/, "")
          const font = await loadFonts(fontAlias, this.staticInfo.url + `fonts/${curFont}`)
          iframeDoc.fonts.add(font) // document.fonts.add(font)
          this.updateBSTableColumnStyle(curColumn, {"font-family": fontAlias})
        }
      }

      const divBGColor = document.createElement("div")
      if ("Background Color") {
        const oldFontSize = queryCellStyle(columnStyle, "background-color", "#ffffff")
        // divBGColor.onclick = () => {}
        divBGColor.className = "mt-5 row"
        divBGColor.innerHTML = `<input type="color" class="form-control form-control-color" value="${oldFontSize}" title="Font color">`
        const inputColor = divBGColor.querySelector("input")
        inputColor.style["max-width"] = "5em"
        inputColor.onchange = (e) => {
          const inputColorValue = e.target.value
          this.updateBSTableColumnStyle(curColumn, {"background-color": inputColorValue})
        }
      }


      const NewToggleBtnEvent = (divNode, btnLabelName, target, attrName, subFunc = () => {
      }) => {
        divNode.className = "mt-5 row"
        divNode.innerHTML = `
<label class="ps-0">${btnLabelName}</label>
<label class="switch">
  <input type="checkbox"><span class="slider round"></span>
</label>`
        const checkSortable = divNode.querySelector(`input[type="checkbox"]`)
        checkSortable.checked = target[attrName] // previous state
        checkSortable.onclick = (e) => {
          target[attrName] = e.target.checked
          subFunc()
          this.updateBSTableColumnStyle(curColumn)
        }
      }

      const NewSliderRangeEvent = (divNode, labelName, min, max, target, attrName, subFunc = () => {
      }) => {
        const regex = new RegExp("[^\\d]+$", 'gm') // find unit px, %, em,...
        const targetValue = target[attrName]
        const matchArray = String(targetValue).match(regex)
        const unit = matchArray !== null ? matchArray[0] : ""
        divNode.className = "mt-2"
        divNode.innerHTML = `
<label class="ps-0">${labelName}(${targetValue})</label>
<input type="range" min="${min}" max="${max}" value="${String(targetValue).replace(regex, "")}">` // // typeof targetValue === "string" ? targetValue.replace('%', "") : targetValue


        const label = divNode.querySelector(`label`)
        const rangeInput = divNode.querySelector(`input[type="range"]`)
        rangeInput.oninput = (e) => {
          const val = e.target.value
          label.innerText = `${labelName}(${val}${unit})`
          target[attrName] = val
          subFunc()
          this.updateBSTableColumnStyle(curColumn)
        }
      }

      const NewFieldSet = (legendName) => {
        const fieldset = document.createElement("fieldset")
        fieldset.className = "mt-2 row"
        fieldset.innerHTML = `
<fieldset>
<legend>${legendName}</legend>
</fieldset>
`
        return fieldset
      }

      const fieldsetImg = NewFieldSet(i18n.LabelImage)
      fieldsetImg.disabled = !(curColumn.imgInfo["isImg"])

      const divResizeImgWidth = document.createElement("div")
      fieldsetImg.append(divResizeImgWidth)
      NewSliderRangeEvent(divResizeImgWidth, i18n.LabelWidth, 1, 300, curColumn.imgInfo, "width")

      const divSortable = document.createElement("div")
      const divIsImg = document.createElement("div")
      NewToggleBtnEvent(divSortable, i18n.LabelSortable, curColumn, "sortable")
      NewToggleBtnEvent(divIsImg, i18n.LabelIsImage, curColumn.imgInfo, "isImg", () => {
        fieldsetImg.disabled = !(curColumn.imgInfo["isImg"])
      })

      const fieldsetColumn = NewFieldSet(i18n.LabelColumnAttr)
      const divColumnWidth = document.createElement("div")
      fieldsetColumn.append(divColumnWidth)
      NewSliderRangeEvent(divColumnWidth, i18n.LabelWidth, 5, 20, curColumn, "width")


      const NewHorizontalRadioButtons = (divNode, labelName, inputArray) => {
        // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/radio#defining_a_radio_group
        divNode.innerHTML = `<label>${labelName}</label>`
        const innerDiv = document.createElement("div")
        for (const value of inputArray) {
          const label = document.createElement("label")
          innerDiv.append(label)
          label.outerHTML = `
<label class="ps-5">
  <input name="${labelName}" type="radio" value="${value}">${value}
</label>`
          const input = innerDiv.querySelector(`input[value="${value}"]`)

          input.onchange = (e) => {
            const val = e.target.value
            curColumn.textAlign = val
            this.updateBSTableColumnStyle(curColumn, {"text-align": val})
          }
        }
        innerDiv.querySelector(`input[value="${curColumn.textAlign}"]`).checked = true // previous state
        divNode.append(innerDiv)
      }

      const divTextAlign = document.createElement("div")
      divTextAlign.className = "mt-3"
      NewHorizontalRadioButtons(divTextAlign, i18n.LabelTextAlign, ["start", "center", "end"])
      fieldsetColumn.append(divTextAlign)

      // combine
      modalBody.append(
        divFont, divBGColor,
        divSortable, divIsImg,
        fieldsetImg, fieldsetColumn
      )
      modal.style.display = "block"
    }

    const iframeCtxWinDoc = this.iframe.contentWindow.document
    const filterID = "selectGroup"
    const filterGroup = iframeCtxWinDoc.getElementById(filterID)
    const optgroup = document.createElement("optgroup")
    optgroup.label = fieldName
    filterGroup.append(optgroup)
    const filterOption = this.options.filter
    if (filterOption.enable && filterOption.isNeedInit) { // must create optgroup first.
      filterGroup.display = "initial"
      const observer = new MutationObserver((mutationRecordList, observer) => {
        for (const mutation of mutationRecordList) {
          switch (mutation.type) {
            case "childList":
              const dataArray = this.table.bootstrapTable('getData', {
                useCurrentPage: false, // all page
                includeHiddenRows: true,
                unfiltered: true, // include all data (unfiltered).
                // formatted:
              })
              for (const curCol of this.columns) {
                if (curCol.isControl) {
                  continue
                }
                const optgroup = iframeCtxWinDoc.querySelector(`optgroup[label="${curCol.field}"]`)
                const colDataArray = dataArray.map(row => {
                  return row[curCol.field]
                })
                const colDataSet = new Set(colDataArray)
                let count = 0
                for (const item of colDataSet) {
                  const option = document.createElement("option")
                  optgroup.append(option)
                  option.value = item
                  option.innerText = item
                  if (++count >= 100) {
                    console.log(`[filter]: Items too much. ${curCol.field}, length:${colDataSet.size} >= 100`)
                    break
                  }
                }
                filterGroup.append(optgroup)
              }
              this.iframe.contentWindow.initMultipleSelect(filterID)
              observer.disconnect()
              return
          }
        }
      })
      observer.observe(this.iframe.contentWindow.document.getElementById("bs-table"), {
        childList: true
      })

      filterGroup.onchange = (event)=>{
        this.updateFilter()
      }

      filterOption.isNeedInit = false
    }

    return `<span>${titleName}<sup> <i class="fas fa-tools" onclick="showPopupConfig(this, '${fieldName},${titleName}')"></i></sup></span>`
  }

  initIframeEvent() {
    this.iframe.onload = (event) => {
      const iframeCtxWin = this.iframe.contentWindow // ??? contentWindow must wait until the HTML is  loaded before it will work.
      const iframeCtxWinDoc = iframeCtxWin.document
      // const TABLE_ID = iframeCtxWinDoc.querySelector("table").getAttribute("id")
      // const TABLE = iframeCtxWinDoc.getElementById(TABLE_ID)
      // this.TABLE = iframeCtxWin.TABLE
      this.table = iframeCtxWin.table // iframeCtxWin.[my-variable] // $("#iframe-bootstrap-table").contents().find(`[id=${TABLE_ID}]`) bootstrap-table???????????????table????????????????????????
      this.table.uniqueId = UNIQUE_ID

      if (this.dataArray.length === 0) {
        return
      }
      const headers = []
      const firstObj = this.dataArray[0]
      for (const [key, value] of Object.entries(firstObj)) {
        headers.push(key)
      }
      this.dataArray = this.dataArray.map((obj, idx) => (obj[UNIQUE_ID] = idx, obj)) // add serial number

      // [refresh bs-table](https://github.com/wenzhixin/bootstrap-table/issues/64)
      const initColumn = (headerName) => {
        return {
          field: headerName,
          // sortable: "true",
          width: 10,
          widthUnit: "%",
          title: this.setConfigColumn(headerName, headerName),
          editable: {
            type: 'text',
            title: headerName,
            emptytext: "",
            validate: function (v) {
              if (!v) return "You can't set the null value on this column"
            },
          },
          // ???? The attribute below is I created not provided by bootstrap-table.
          imgInfo: {
            isImg: false,
            width: "50%",
            height: "auto",
          },
          textAlign: "start",
          isControl: false,
        }
      }
      const columns = headers.map(headerName => (initColumn(headerName)))
      columns.splice(0, 0, // Add a checkbox to select the whole row.
        {
          checkbox: true, width: 2, widthUnit: "%", align: 'center',
          isControl: true
        },
        // {field: UNIQUE_ID, title: "uid", visiable: false, sortable: true} // new column for UNIQUE_ID // It's ok if you aren't gonna show it to the user. Data still save on the datatable, no matter you set the filed or not.
      )
      columns.push({
        field: "Action", title: "Action", align: "center", width: 5, widthUnit: "%",
        isControl: true,
        formatter: (value, row, index, field) => {
          const curID = row[UNIQUE_ID]
          return [
            `<button type="button" class="btn btn-default btn-sm" onclick="DeleteItem(${curID})">`,
            `<span class="glyphicon glyphicon-trash"></span>`,
            `</button>`
          ].join('')
        }
      })

      this.columns = columns
      this.initToolbar()
      this.table.bootstrapTable('refreshOptions',
        {
          columns: columns, // [{field: "Name", title: "??????"}, {field: "Desc", title: "??????"}]
          // url: dataURL // You can use ``url`` instead of ``data``, but there is unnecessary since we already get all the data.
          data: this.dataArray,
          height: 628,
          uniqueId: UNIQUE_ID, // Using the ``headers[0]`` is not a great idea, so I create a column(__id__) instead of it. // data-unique-id
          clickToSelect: true,
          buttonsClass: BUTTONS_CLASS,
          pageSize: 25,
          paginationUseIntermediate: true,
          paginationPagesBySide: 2,
          locale: this.locale,
        }
      )
      // table.uniqueId = UNIQUE_ID // Add other attributes for bootstrap-table
      this.table.bootstrapTable('refresh')
      this.table.highlight($(".search.bs.table").val()) // jquery.highlight-5.min.js
      this.table.on('search.bs.table', (e, text) => {
        this.table.highlight(text)
      })
    }
  }

  initToolbar() {
    const divToolbar = this.iframe.contentWindow.document.getElementById("toolbar")

    let button
    let icon

    const buttonsObj = {
      btnDeleteSelect: {
        text: 'Delete',
        icon: 'fa-trash-alt',
        attributes: {
          style: 'color:#ff0000',
          title: 'Delete all selection data'
        },
        event: {
          'click': () => {
            if (!confirm(i18n.AskDeleteSelection)) {
              return
            }

            const selectDataArray = this.table.bootstrapTable('getSelections')
            const ids = []
            selectDataArray.forEach(obj => {
                ids.push(obj[UNIQUE_ID])
              }
            )
            this.table.bootstrapTable('remove', {
              field: UNIQUE_ID,
              values: ids,
            })
          }
        }
      },
      btnUsersAdd: {
        text: 'Save',
        icon: 'fa-save',
        attributes: {
          title: 'Save data'
        },
        event: {
          'click': () => {
            alert("save data to server")
          },
          'mouseenter': () => {
          },
          'mouseleave': () => {
          }
        },
      }
    }

    for (const [key, btnObj] of Object.entries(buttonsObj)) {
      const iconClass = btnObj.icon
      button = document.createElement("button")
      button.setAttribute("class", `btn btn-${BUTTONS_CLASS}`)
      button.setAttribute("style", btnObj.attributes.style)
      button.setAttribute("title", btnObj.attributes.title)
      button.addEventListener("click", btnObj.event.click)
      icon = document.createElement("i")
      icon.setAttribute("class", `far ${iconClass}`)
      button.appendChild(icon)
      divToolbar.appendChild(button)
    }
  }
}

/**
 * @param {StaticInfo} staticInfo The properties.
 * @param {string} lang
 */
function inputFileHandler(staticInfo, lang) {
  const inputFile = document.getElementById("uploadFile")
  const inputValue = inputFile.value
  if (inputValue === "") {
    return null
  }

  const selectedFile = document.getElementById('uploadFile').files[0] // https://developer.mozilla.org/en-US/docs/Web/API/File/Using_files_from_web_applications
  const promise = new Promise(resolve => {
    const fileContent = ReadFile(selectedFile)
    resolve(fileContent)
  })

  promise.then(fileContent => {
    fileContent = removeExtraSpace(fileContent)
    const options = {
      separator: ",", // "\t",
      delimiter: '"', // default "
      headers: true // default true
    }
    const dataArray = $.csv.toObjects(fileContent, options) // jquery.csv.min.js
    const bsTable = new BSTable(dataArray, staticInfo, lang)
    return null
  })
}


async function inputStaticDirHandler() {
  const inputStaticDir = document.getElementById("inputStaticDir")
  const inputStaticDirValue = inputStaticDir.value
  if (inputStaticDirValue === "") {
    return new FileInfo()
  }

  const formData = new FormData()
  formData.set("inputStaticDir", inputStaticDirValue)

  // https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
  const response = await fetch("/api/os/Stat", {
    method: "post",
    body: formData
  })

  if (!response.ok) {
    const errMsg = await response.text()
    throw Error(`${response.statusText} (${response.status}) | ${errMsg} `)
  }
  const obj = await response.json()
  return new FileInfo(obj.Name, obj.Size, obj.ModTime, obj.IsDir, obj.Path)
}

/**
 * @param {FileInfo} fileInfo
 * @return {string} static URL
 */
async function AskInitStaticDir(fileInfo) {
  if (fileInfo.Path === undefined)
    return ""
  const formData = new FormData()
  formData.set("staticInfoObj", JSON.stringify(fileInfo))

  const response = await fetch("/user/static/", {
    method: "post",
    body: formData
  })

  if (!response.ok) {
    const errMsg = await response.text()
    throw Error(`${response.statusText} (${response.status}) | ${errMsg} `)
  }
  const obj = await response.json()
  return obj.staticDirURL
}

async function getLocale() {
  const response = await fetch("/config/", {method: "get"})
  if (!response.ok) {
    return "en-us"
  }
  const obj = await response.json()
  return obj.Lang
}

async function onCommit() {
  document.getElementById("input&commit").open = false
  let staticFileInfo = new FileInfo()
  let staticInfo = new StaticInfo("", staticFileInfo)
  new Promise((resolve, reject) => {
    resolve(inputStaticDirHandler())
  })
    .then((staticInfoObj) => {
      staticFileInfo = staticInfoObj
      return AskInitStaticDir(staticInfoObj)
    })
    .then(staticURL => {
      staticInfo = new StaticInfo(staticURL, staticFileInfo)
      return ""
    })
    .then((msg) => {
      return getLocale()
    })
    .then((lang) => {
      inputFileHandler(staticInfo, lang)
    })
    .catch((error) => {
      alert(error)
    })
}

(
  () => {
    window.onload = () => {
      // document.getElementById("uploadFile").onchange = () => {}
      const commitBtn = document.getElementById("Commit")
      const aboutBtn = document.getElementById("About")
      commitBtn.onclick = () => {
        // https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Asynchronous/Async_await
        // https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
        onCommit().then()
      }
      aboutBtn.onclick = () => {
        window.location.href = "/about/"
      }
      document.getElementById("ExitApp").onclick = () => {
        window.location.href = "/shutdown/"
      }
    }
  }
)()


