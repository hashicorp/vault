export default function(fileName, toTrimArray = []) {
  const extension = new RegExp(toTrimArray.join('$|'));
  return fileName.replace(extension, '');
}
