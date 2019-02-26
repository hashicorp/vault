import RouteRecognizer from 'route-recognizer';

const {
  Normalizer: { normalizePath, encodePathSegment },
} = RouteRecognizer;

export function encodePath(path) {
  return path
    .split('/')
    .map(encodePathSegment)
    .join('/');
}

export { normalizePath, encodePathSegment };
