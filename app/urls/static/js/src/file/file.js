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

class BSTable {
  constructor(dataArray) {
    this.dataArray = dataArray
    const div = document.getElementById('div-csv-data')
    div.querySelectorAll('*').forEach(node => node.remove())
    this.table = undefined
    this.iframe = this.newIframe()
    this.initIframeEvent()
    div.append(this.iframe)
  }

  newIframe() {
    const iframe = document.createElement('iframe')
    iframe.id = IFRAME_ID
    iframe.height = "100%"
    iframe.width = "100%"
    iframe.src = "/bs-table/"
    return iframe
  }

  initIframeEvent() {
    this.iframe.onload = (event) => {
      const iframeCtxWin = this.iframe.contentWindow // ★ contentWindow must wait until the HTML is  loaded before it will work.
      const iframeCtxWinDoc = iframeCtxWin.document
      // const TABLE_ID = iframeCtxWinDoc.querySelector("table").getAttribute("id")
      // const TABLE = iframeCtxWinDoc.getElementById(TABLE_ID)
      // this.TABLE = iframeCtxWin.TABLE
      this.table = iframeCtxWin.table // iframeCtxWin.[my-variable] // $("#iframe-bootstrap-table").contents().find(`[id=${TABLE_ID}]`) bootstrap-table導入的時候table就會消失會取不到
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
      const columns = headers.map(e => ({
        field: e, title: e, sortable: "true",
        editable: {
          type: 'text',
          title: e,
          emptytext: "",
          validate: function (v) {
            if (!v) return "You can't set the null value on this column"
          },
        }
      }))
      columns.splice(0, 0,
        {checkbox: true, width: 64, align: 'center'}, // Add a checkbox to select the whole row.
        // {field: UNIQUE_ID, title: "uid", visiable: false, sortable: true} // new column for UNIQUE_ID // It's ok if you aren't gonna show it to the user. Data still save on the datatable, no matter you set the filed or not.
      )
      columns.push({
        field: "tableAction", title: "Action", align: "center", width: 64,
        formatter: (value, row, index, field) => {
          const curID = row[UNIQUE_ID]
          return [
            `<button type="button" class="btn btn-default btn-sm" onclick="DeleteItem(${curID})">`,
            `<span class="glyphicon glyphicon-trash"></span>`,
            `</button>`
          ].join('')
        }
      })
      this.initToolbar()
      this.table.bootstrapTable('refreshOptions',
        {
          columns: columns, // [{field: "Name", title: "名稱"}, {field: "Desc", title: "說明"}]
          // url: dataURL // You can use ``url`` instead of ``data``, but there is unnecessary since we already get all the data.
          data: this.dataArray,
          height: 768,
          uniqueId: UNIQUE_ID, // Using the ``headers[0]`` is not a great idea, so I create a column(__id__) instead of it. // data-unique-id
          clickToSelect: true,
          buttonsClass: BUTTONS_CLASS,
        }
      )
      // table.uniqueId = UNIQUE_ID // Add other attributes for bootstrap-table
      this.table.bootstrapTable('refresh')
    }
  }

  initToolbar() {
    const divToolbar = this.iframe.contentWindow.document.getElementById("toolbar")
    divToolbar.setAttribute("style", "padding-left:0.4em")
    divToolbar.setAttribute("class", "columns columns-right btn-group float-right")
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
            if (!confirm('Are you sure you want to delete all selection data?')) {
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

(
  () => {
    window.onload = () => {
      const inputFile = document.getElementById("uploadFile")
      inputFile.onchange = () => {
        const inputValue = inputFile.value
        if (inputValue === "") {
          return
        }

        const selectedFile = document.getElementById('uploadFile').files[0] // https://developer.mozilla.org/en-US/docs/Web/API/File/Using_files_from_web_applications
        const promise = new Promise(resolve => {
          const fileContent = ReadFile(selectedFile)
          resolve(fileContent)
        })

        promise.then(fileContent => {
          fileContent = removeExtraSpace(fileContent)
          const dataArray = $.csv.toObjects(fileContent)
          const bsTable = new BSTable(dataArray)
        })
      }
    }
  }
)()


