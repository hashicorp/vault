/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export function keyIsFolder(key: string) {
  return key ? !!key.match(/\/$/) : false;
}

export function keyPartsForKey(key: string) {
  if (!key) {
    return null;
  }
  const isFolder = keyIsFolder(key);
  const parts = key.split('/');
  if (isFolder) {
    // remove last item which is empty
    parts.pop();
  }
  return parts.length > 1 ? parts : null;
}

export function parentKeyForKey(key: string) {
  const parts = keyPartsForKey(key);
  if (!parts) {
    return '';
  }
  return parts.slice(0, -1).join('/') + '/';
}

export function keyWithoutParentKey(key: string) {
  return key ? key.replace(parentKeyForKey(key), '') : null;
}

export function ancestorKeysForKey(key: string) {
  const ancestors = [];
  let parentKey = parentKeyForKey(key);

  while (parentKey) {
    ancestors.unshift(parentKey);
    parentKey = parentKeyForKey(parentKey);
  }

  return ancestors;
}
