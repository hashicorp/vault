/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export function pathIsDirectory(pathToSecret) {
  // This regex only checks for / at the end of the string. ex: boop/ === true, boop/bop === false;
  return pathToSecret ? !!pathToSecret.match(/\/$/) : false;
}

export function pathIsFromDirectory(path) {
  // This regex just looks for a / anywhere in the path. ex: boop/ === true, boop/bop === true;
  return path ? !!path.match(/\//) : false;
}

function splitSegments(secretPath) {
  const segments = secretPath.split('/').filter((path) => path);
  segments.map((_, index) => {
    return segments.slice(0, index + 1).join('/');
  });
  return segments.map((segment, idx) => {
    return {
      label: segment,
      model: segments.slice(0, idx + 1).join('/'),
    };
  });
}

/**
 * breadcrumbsForSecret is for generating page breadcrumbs for a secret path
 * @param {string} secretPath is the full path to secret (like 'my-secret' or 'beep/boop')
 * @param {boolean} lastItemCurrent
 * @returns array of breadcrumbs specific to KV engine
 */
export function breadcrumbsForSecret(backend, secretPath, lastItemCurrent = false) {
  if (!backend || !secretPath) return [];
  const isDir = pathIsDirectory(secretPath);
  const segments = splitSegments(secretPath);

  return segments.map((segment, index) => {
    if (index === segments.length - 1) {
      if (lastItemCurrent) {
        return {
          label: segment.label,
        };
      }
      if (!isDir) {
        return { label: segment.label, route: 'secret.index', models: [backend, segment.model] };
      }
    }
    return { label: segment.label, route: 'list-directory', models: [backend, `${segment.model}/`] };
  });
}
