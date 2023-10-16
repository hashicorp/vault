export function sanitizePath(path) {
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}

export function ensureTrailingSlash(path) {
  return path.replace(/(\w+[^/]$)/g, '$1/');
}
