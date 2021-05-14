import * as base64 from "../../pkg/encoding/base64.js"
import {ReadFile} from "../../pkg/io/ioutil.js"

const IFRAME_ID = "iframe-bootstrap-table"
const UNIQUE_ID = "__id__"

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

function createBSTable(dataArray) {
  const div = document.getElementById('div-csv-data')
  div.querySelectorAll('*').forEach(node => node.remove())
  const iframe = document.createElement('iframe')
  iframe.id = IFRAME_ID
  iframe.height = "100%"
  iframe.width = "100%"
  iframe.src = "/bs-table/"
  div.append(iframe)
  // iframe.contentWindow.location.reload(false)

  iframe.onload = (event) => {
    const iframeCtxWin = iframe.contentWindow
    const iframeCtxWinDoc = iframeCtxWin.document

    // const TABLE_ID = iframeCtxWinDoc.querySelector("table").getAttribute("id")
    // const TABLE = iframeCtxWinDoc.getElementById(TABLE_ID)
    const TABLE = iframeCtxWin.TABLE
    const table = iframeCtxWin.table // iframeCtxWin.[my-variable] // $("#iframe-bootstrap-table").contents().find(`[id=${TABLE_ID}]`) bootstrap-table導入的時候table就會消失會取不到

    if (dataArray.length === 0) {
      return
    }
    const headers = []
    const firstObj = dataArray[0]
    for (const [key, value] of Object.entries(firstObj)) {
      headers.push(key)
    }
    dataArray = dataArray.map((obj, idx) => (obj[UNIQUE_ID] = idx, obj)) // add serial number

    // [refresh bs-table](https://github.com/wenzhixin/bootstrap-table/issues/64)
    const columns = headers.map(e => ({field: e, title: e, sortable: "true"}))
    columns.splice(0, 0,
      {checkbox: true, width: 64, align: 'center'}, // Add a checkbox to select the whole row.
    )
    table.bootstrapTable('refreshOptions',
      {
        columns: columns, // [{field: "Name", title: "名稱"}, {field: "Desc", title: "說明"}]
        // url: dataURL // You can use ``url`` instead of ``data``, but there is unnecessary since we already get all the data.
        data: dataArray,
        height: 768,
        uniqueId: UNIQUE_ID, // Using the ``headers[0]`` is not a great idea, so I create a column(__id__) instead of it. // data-unique-id
        buttons: { // data-buttons
          btnDeleteSelect: {
            text: 'Delete',
            icon: 'fa-trash-alt',
            attributes: {
              style: 'color:#82c91e',
              title: 'Delete all selection data'
            },
            event: {
              'click': () => {
                dataArray = table.bootstrapTable('getSelections')
                const ids = []
                dataArray.forEach(obj => {
                    ids.push(obj[UNIQUE_ID])
                  }
                )
                table.bootstrapTable('remove', {
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
                save2Server()
              },
              'mouseenter': () => {
              },
              'mouseleave': () => {
              }
            },
          }
        }
      }
    )
    // table.uniqueId = UNIQUE_ID // Add other attributes for bootstrap-table
    table.bootstrapTable('refresh')
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
          createBSTable(dataArray)
        })
      }
    }
  }
)()


