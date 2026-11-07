/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

function encodePath(path) {
  return path
    ? path
        .split('/')
        .map((segment) => encodeURIComponent(segment))
        .join('/')
    : path;
}

function normalizePath(path) {
  // Unlike normalizePath from route-recognizer, this method assumes
  // we do not have percent-encoded data octets as defined in
  // https://datatracker.ietf.org/doc/html/rfc3986
  return path
    ? path
        .split('/')
        .map((segment) => decodeURIComponent(segment))
        .join('/')
    : '';
}

function recoverEngineMountPath(backend) {
  return backend?.includes('/') ? backend.split('/').join('%2F') : backend || '';
}

export { normalizePath, encodePath, recoverEngineMountPath };
