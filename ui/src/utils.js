import convert from "xml-js";

export function xmlToTree(xml) {
    let tree = convert.xml2js(xml, {
        compact: false,
        ignoreComment: true,
        elementsKey: 'children',
    }).children[0].children
    // walk tree
    walkNode(tree, 1)
    return tree
}

export function isIterable(object) {
    return object != null && typeof object[Symbol.iterator] === 'function'
}

/**
 * @param tree
 * @param counter{int}
 * @return {int}
 */
function walkNode(tree, counter) {
    // Parent first traverse
    for (let node of tree) node.id = counter++
    for (let node of tree) {
        if (node.children) counter = walkNode(node.children, counter)
    }
    return counter
}

/**
 * @param millis{number}
 * @return {Promise<any>}
 */
export async function sleep(millis) {
    return new Promise(resolve => setTimeout(resolve, millis))
}
