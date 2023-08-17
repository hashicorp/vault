/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export function pathIsDirectory(pathToSecret) {
  // This regex only checks for / at the end of the string. ex: boop/ === true, boop/bop === false;
  return pathToSecret ? !!pathToSecret.match(/\/$/) : false;
}

export function pathIsFromDirectory(path) {
  // This regex just looks for a / anywhere in the path. ex: boop/ === true, boop/bop === true;
  return path ? !!path.match(/\//) : false;
}

function finalCrumb(segment, isLast, modelPath) {
  const model = modelPath || segment;
  if (isLast) {
    return { label: segment };
  }
  return { label: segment, route: 'secret.details', model };
}

/**
 *
 * @param {string} secretPath is the full path to secret (like 'my-secret' or 'beep/boop')
 * @param {boolean} lastItemCurrent
 * @param {*} options
 * @returns
 */
export function breadcrumbsForSecret(secretPath, lastItemCurrent = false) {
  const pathAsArray = secretPath.split('/');
  const modelIdArray = pathAsArray.map((_, index) => {
    return pathAsArray.slice(0, index + 1).join('/');
  });

  return pathAsArray
    .map((key, index) => {
      if (!key) {
        // path segment is empty which means path is a directory, return null so it can be filtered out
        return null;
      }
      if (pathAsArray.length - 1 === index) {
        return finalCrumb(key, lastItemCurrent, modelIdArray[index]);
      }
      return { label: key, route: 'list-directory', model: `${modelIdArray[index]}/` };
    })
    .filter((segment) => segment);
}
