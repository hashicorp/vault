// will trim a given set of endings from the end of a string
// if isExtension is true, the first char of that string will be escaped
// in the regex
export default function(str, endings = [], isExtension = true) {
  let trimRegex = new RegExp(endings.map(ext => `${ext}$`).join('|'));
  if (isExtension) {
    trimRegex = new RegExp(endings.map(ext => `\\${ext}$`).join('|'));
  }
  return str.replace(trimRegex, '');
}
