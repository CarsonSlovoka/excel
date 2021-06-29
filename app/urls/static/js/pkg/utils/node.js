/**
 * Create a new Node
 *
 * @param node {string}
 * @param attrObj {Object}
 * @returns {string}
 * @example
 *  const myDemoDiv document.querySelector("div.test[data-test='demo']")
 *  const svg = GetNode("svg", {id: "my-svg", style: "width:60vw; height:40vh;"})
 *  myDemoDiv.replaceWith(svg)
 * @see https://stackoverflow.com/a/37411738
 *  - https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelector
 */
export function GetNode(node, attrObj) {
  node = document.createElementNS("http://www.w3.org/2000/svg", node);
  for (let attr in attrObj)
    node.setAttributeNS(null, attr.replace(/[A-Z]/g, function(m, p, o, s) { return "-" + m.toLowerCase(); }), attrObj[attr]);
  return node
}
