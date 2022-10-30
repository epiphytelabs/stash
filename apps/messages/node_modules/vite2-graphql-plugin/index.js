const os = require('os')
const gql = require('graphql-tag')

const generateOutput = require('./utils/generateOutput')
const expandImports = require('./utils/expandImpots')

const transform = (code) => {
  const doc = gql`
    ${code}
  `
  let headerCode = `var doc = ${JSON.stringify(
    doc
  )}; doc.loc.source = ${JSON.stringify(doc.loc.source)};`

  let outputCode = ''

  // Allow multiple query/mutation definitions in a file. This parses out dependencies
  // at compile time, and then uses those at load time to create minimal query documents
  // We cannot do the latter at compile time due to how the #import code works.
  const operationCount = doc.definitions.reduce((accum, op) => {
    return op.kind === 'OperationDefinition' ? accum + 1 : accum
  }, 0)

  if (operationCount < 1) {
    outputCode += '\nexport default doc;'
  } else {
    outputCode += generateOutput()

    for (const op of doc.definitions) {
      if (op.kind === 'OperationDefinition') {
        if (!op.name) {
          if (operationCount > 1) {
            throw 'Query/mutation names are required for a document with multiple definitions'
          } else {
            continue
          }
        }

        const opName = op.name.value
        outputCode += `export const ${opName} = oneQuery(doc, "${opName}");`
      }
    }
  }

  const importOutputCode = expandImports(code, doc)

  return headerCode + os.EOL + importOutputCode + os.EOL + outputCode + os.EOL
}

const fileRegex = /\.(graphql)$/

module.exports = function viteGraphQlPluing() {
  return {
    name: 'vite-plugin-graphql',
    async transform(src, id) {
      if (fileRegex.test(id)) {
        return {
          code: transform(src),
          map: null // provide source map if available
        }
      }
    }
  }
}
