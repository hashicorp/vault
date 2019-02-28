import RouteRecognizer from 'route-recognizer';

const {
  Normalizer: { normalizePath, encodePathSegment },
} = RouteRecognizer;

export function encodePath(path) {
  return path
    ? path
        .split('/')
        .map(encodePathSegment)
        .join('/')
    : path;
}

export { normalizePath, encodePathSegment };
