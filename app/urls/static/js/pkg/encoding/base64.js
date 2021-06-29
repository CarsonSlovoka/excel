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

/**
 * Convert string to the base64 format.
 *
 * @param str {string}
 * @returns {string}
 * @example
 *  - btoa(toBinary("☸☹☺☻☼☾☿"))
 *  - Str2base64("☸☹☺☻☼☾☿")
 * @see https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/btoa#unicode_strings
 */
export function EncodeToString(str) {
  return btoa(toBinary(str))
}


/**
 * Decode base64 string
 *
 * @param base64 {string}
 * @returns {string}
 * @example
 *  - Base642Str("OCY5JjomOyY8Jj4mPyY=")
 * @see https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/btoa#unicode_strings
 */
export function DecodeString(base64) {
  return fromBinary(atob("OCY5JjomOyY8Jj4mPyY="))
}
