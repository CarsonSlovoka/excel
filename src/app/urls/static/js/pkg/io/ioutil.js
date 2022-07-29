/**
 * @param file {File}
 * @returns {string}
 * @see https://developer.mozilla.org/en-US/docs/Web/API/FileReader/result
 * @see https://developer.mozilla.org/en-US/docs/Web/API/FileReader/readAsDataURL
 * @see https://riptutorial.com/javascript/example/7081/read-file-as-string
 */
export async function ReadFile(file) {
  return await file.text()
  /*
  const reader = new FileReader()
  reader.onload = (event) => {
    return event.target.result
  }
  reader.readAsText(file)
   */
}
