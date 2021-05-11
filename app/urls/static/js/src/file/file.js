import * as base64 from "../../pkg/encoding/base64.js"
import * as node from "../../pkg/utils/node.js"

function createBSTable(dataURL) {
  const div = document.getElementById('div-csv-data')
  div.querySelectorAll('*').forEach(node => node.remove())
  const iframe = document.createElement('iframe')
  iframe.id = "iframe-bootstrap-table"
  iframe.height = "100%"
  iframe.width = "100%"
  dataURL = document.baseURI + dataURL // 因為我們有改變iframe.src，所以相對路徑會變成/bs-table/ 這並不是我們所期望的，所以要先決定URL的絕對路徑
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

    // TABLE.setAttribute("data-unique-id", "id")
    // TABLE.setAttribute("data-buttons", "buttons")

    const xhr = new XMLHttpRequest()
    xhr.open("GET", dataURL, true)
    xhr.onload = (progressEvent) => {
      if (xhr.status !== 200) {
        alert(`${xhr.responseText}`)
        return
      }

      // [Parsing JSON from XmlHttpRequest.responseJSON](https://stackoverflow.com/a/8416963)
      const dataArray = JSON.parse(xhr.responseText) // array
      if (dataArray.length === 0) {
        return
      }
      const headers = []
      // dataArray.forEach(rowObj => {})
      const firstObj = dataArray[0]
      for (const [key, value] of Object.entries(firstObj)) {
        headers.push(key)
      }

      // [refresh bs-table](https://github.com/wenzhixin/bootstrap-table/issues/64)
      const columns = headers.map(e => ({field: e, title: e}))
      table.bootstrapTable('refreshOptions',
        {
          columns: columns, // [{field: "Name", title: "名稱"}, {field: "Desc", title: "說明"}]
          // url: dataURL // You can use ``url`` instead of ``data``, but there is unnecessary since we already get all the data.
          data: dataArray
        }
      )
      table.bootstrapTable('refresh')
    }
    // xhr.responseType = "json"
    xhr.send()
  }
}

(
  () => {
    const myDemoDiv = document.querySelector("div.test[data-test='demo']")
    const svg = node.GetNode("svg", {id: "my-svg", style: "width:60vw; height:40vh;"})
    myDemoDiv.replaceWith(svg)
    window.onload = () => {
      const inputFile = document.getElementById("uploadFile")
      inputFile.onchange = () => {
        const inputValue = inputFile.value
        if (inputValue === "") {
          return
        }
        const filename = inputValue.substring(inputValue.lastIndexOf('\\') + 1).split(".")[0] // get basename without extension
        const filenameB64 = base64.Str2base64(filename) // Ensure Chinese filename is working.

        const xhr = new XMLHttpRequest()
        xhr.open("POST", filenameB64, true)
        // xhr.setRequestHeader("Content-Type", "multipart/form-data; charset=utf-8") // do not change the content type Otherwise will occurring error of the "no multipart boundary param in Content-Type."

        xhr.onload = (oEvent) => {
          if (xhr.status === 200) {
            const targetCSVURL = filenameB64 + "/data" // /{worksheets}/data
            createBSTable(targetCSVURL)
          } else {
            alert(`${xhr.responseText}`)
          }
        }

        const selectedFile = document.getElementById('uploadFile').files[0] // https://developer.mozilla.org/en-US/docs/Web/API/File/Using_files_from_web_applications
        const formData = new FormData() // https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects
        formData.set("myUploadFile", selectedFile)
        xhr.send(formData)
      }
    }
  }
)()
