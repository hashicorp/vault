function keyIsFolder(key) {
  return key ? !!key.match(/\/$/) : false;
}

function keyPartsForKey(key) {
  if (!key) {
    return null;
  }
  var isFolder = keyIsFolder(key);
  var parts = key.split('/');
  if (isFolder) {
    parts.pop();
  }
  return parts.length > 1 ? parts : null;
}

function parentKeyForKey(key) {
  var parts = keyPartsForKey(key);
  if (!parts) {
    return null;
  }
  return parts.slice(0, -1).join('/') + '/';
}

function keyWithoutParentKey(key) {
  return key ? key.replace(parentKeyForKey(key), '') : null;
}

function ancestorKeysForKey(key) {
  var ancestors = [],
    parentKey = parentKeyForKey(key);

  while (parentKey) {
    ancestors.unshift(parentKey);
    parentKey = parentKeyForKey(parentKey);
  }

  return ancestors.length ? ancestors : null;
}

export default {
  keyIsFolder,
  keyPartsForKey,
  parentKeyForKey,
  keyWithoutParentKey,
  ancestorKeysForKey,
};
