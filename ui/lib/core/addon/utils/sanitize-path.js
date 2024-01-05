/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export function sanitizePath(path) {
  if (!path) return '';
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}

export function sanitizeStart(path) {
  if (!path) return '';
  //remove leading slashes
  return path.trim().replace(/^\/+/, '');
}

export function ensureTrailingSlash(path) {
  return path.replace(/(\w+[^/]$)/g, '$1/');
}

/**
 * getRelativePath is for removing matching segments of a subpath from the front of a full path.
 * This method assumes that the full path starts with all of the root path.
 * @param {string} fullPath eg apps/prod/app_1/test
 * @param {string} rootPath eg apps/prod
 * @returns the leftover segment, eg app_1/test
 */
export function getRelativePath(fullPath = '', rootPath = '') {
  const root = sanitizePath(rootPath);
  const full = sanitizePath(fullPath);

  if (!root) {
    return full;
  } else if (root === full) {
    return '';
  }
  return sanitizePath(full.substring(root.length));
}
