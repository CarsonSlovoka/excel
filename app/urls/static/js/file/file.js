function toBinary(string) {
  // https://stackoverflow.com/a/67415709
  const codeUnits = new Uint16Array(string.length);
  for (let i = 0; i < codeUnits.length; i++) {
    codeUnits[i] = string.charCodeAt(i);
  }
  return String.fromCharCode(...new Uint8Array(codeUnits.buffer));
}

function fromBinary(binary) {
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return String.fromCharCode(...new Uint16Array(bytes.buffer));
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

        const xhr = new XMLHttpRequest()
        xhr.open("POST", inputValue, true)
        // xhr.setRequestHeader("Content-Type", "multipart/form-data; charset=utf-8") // do not change the content type Otherwise will occurring error of the "no multipart boundary param in Content-Type."

        xhr.onload = function(oEvent) {
          if (xhr.status === 200) {
            alert('save successfully')
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
