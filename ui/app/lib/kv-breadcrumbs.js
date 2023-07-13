/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * This file is based off the older key-utils. However, after moving KV to it's own engine we needed to modify some of the methods here for the new models.
 * */

function pathIsFolder(secretPrefix) {
  // ex: boop/ === true, boop/bop === false; This regex only checks for the end of the string.
  return secretPrefix ? !!secretPrefix.match(/\/$/) : false;
}

function pathIsFromNested(path) {
  // ex: beep/boop/bop doesn't need to end in a / just include on in the path.
  return path ? !!path.match(/\//) : false;
}

function secretPrefixParts(key) {
  if (!key) {
    return null;
  }
  var isFolder = pathIsFolder(key);
  var parts = key.split('/');
  if (isFolder) {
    parts.pop();
  }
  return parts.length > 0 ? parts : null;
}

function modelIdsForSecretPrefix(secretPrefix) {
  const secretPrefixAsArray = secretPrefix.split('/');
  secretPrefixAsArray.pop(); // remove the last / so you can get the correct index count

  const modelIdArray = secretPrefixAsArray.map((key, index) => {
    return `${secretPrefixAsArray.slice(0, index + 1).join('/')}/`;
  });
  const secretPrefixPartsAsAnArray = secretPrefixParts(secretPrefix);

  return secretPrefixPartsAsAnArray.map((key, index) => {
    // we do not want to return "route or model" on the last item otherwise it will add link to the current page.
    if (secretPrefixPartsAsAnArray.length - 1 === index) {
      return { label: key };
    }
    return { label: key, route: 'list-nested-secret', model: modelIdArray[index] };
  });
}

export { modelIdsForSecretPrefix, pathIsFolder, pathIsFromNested };
