/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * This file is based off the older key-utils. However, after moving KV to it's own engine we needed to modify some of the methods here for the new models.
 * */

function secretPrefixIsFolder(secretPrefix) {
  return secretPrefix ? !!secretPrefix.match(/\/$/) : false;
}

function keyPartsForKey(key) {
  if (!key) {
    return null;
  }
  var isFolder = secretPrefixIsFolder(key);
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

  return keyPartsForKey(secretPrefix).map((key, index) => {
    return { label: key, route: 'list-nested-secret', model: modelIdArray[index] };
  });
}

export { modelIdsForSecretPrefix, secretPrefixIsFolder };
