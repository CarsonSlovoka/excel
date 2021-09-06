export class AskInputFile {

  static BuildAll() {
    document.querySelectorAll(`input[data-com="ask-file"]`).forEach(input => {
      new AskInputFile(input)
    })
  }

  constructor(inputNode, labelName) {
    labelName = labelName ?? inputNode.placeholder
    const labelClass = " " + inputNode.dataset.labelClass ?? ""
    const frag = document.createRange().createContextualFragment(`
  <label class='form-label${labelClass}'>${labelName}
    ${inputNode.outerHTML}
  </label>
        `)
    const input = frag.querySelector(`input`)

    inputNode.parentNode.insertBefore(frag, inputNode)
    inputNode.remove()

    input.type = "file"
    input.classList.add("mt-2", "form-control")

    input.ondragleave = dragEvent => {
      dragEvent.target.classList.remove("border", "border-3", "border-primary")
    }

    input.ondragover = input.ondragenter = (dragEvent) => {
      dragEvent.preventDefault()
      dragEvent.target.classList.add("border", "border-3", "border-primary")
    }

    input.ondrop = (dragEvent) => {
      dragEvent.preventDefault()
      input.classList.remove("border", "border-3", "border-primary")

      for (let i = 0; i < dragEvent.dataTransfer.files.length; i++) {
        const file = dragEvent.dataTransfer.files[i]
        if (!input.accept.split(",").includes(file.type)) {
          throw new Error(`Invalid file format: ${file.type}. ${input.accept} expected.`)
        }
      }
      if (input.multiple) {
        input.files = dragEvent.dataTransfer.files
      } else {
        const dataTransfer = new DataTransfer()
        dataTransfer.items.add(dragEvent.dataTransfer.files[0])
        input.files = dataTransfer.files
      }
    }
  }
}

/*
(() => {
  window.addEventListener(`load`, () => AskInputFile.BuildAll(), {once: true})
})()
 */
