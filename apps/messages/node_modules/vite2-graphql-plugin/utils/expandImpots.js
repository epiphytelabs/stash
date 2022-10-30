module.exports = function expandImports(source) {
  const lines = source.split(/\r\n|\r|\n/)
  let outputCode = `
    var names = {};
    function unique(defs) {
      return defs.filter(
        function(def) {
          if (def.kind !== 'FragmentDefinition') return true;
          var name = def.name.value
          if (names[name]) {
            return false;
          } else {
            names[name] = true;
            return true;
          }
        }
      )
    }`

  lines.some((line) => {
    if (line[0] === '#' && line.slice(1).split(' ')[0] === 'import') {
      const importFile = line.slice(1).split(' ')[1]
      const parseDocument = `require(${importFile})`
      const appendDef = `doc.definitions = doc.definitions.concat(unique(${parseDocument}.definitions));`
      outputCode += appendDef + os.EOL
    }
    return line.length !== 0 && line[0] !== '#'
  })

  return outputCode
}
