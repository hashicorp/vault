/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import { FORM } from './kv-selectors';
import { encodePath } from 'vault/utils/path-encoding-helpers';

// CUSTOM ACTIONS RELEVANT TO KV-V2

export const writeSecret = async function (backend, path, key, val, ns = null) {
  const url = `vault/secrets/${backend}/kv/create`;
  ns ? await visit(url + `?namespace=${ns}`) : await visit(url);
  await fillIn(FORM.inputByAttr('path'), path);
  await fillIn(FORM.keyInput(), key);
  await fillIn(FORM.maskedValueInput(), val);
  return click(FORM.saveBtn);
};

export const writeVersionedSecret = async function (backend, path, key, val, version = 2, ns = null) {
  await writeSecret(backend, path, 'key-1', 'val-1', ns);
  for (let currentVersion = 2; currentVersion <= version; currentVersion++) {
    const url = `/vault/secrets/${backend}/kv/${encodeURIComponent(path)}/details/edit`;
    ns ? await visit(url + `?namespace=${ns}`) : await visit(url);

    if (currentVersion === version) {
      await fillIn(FORM.keyInput(), key);
      await fillIn(FORM.maskedValueInput(), val);
    } else {
      await fillIn(FORM.keyInput(), `key-${currentVersion}`);
      await fillIn(FORM.maskedValueInput(), `val-${currentVersion}`);
    }
    await click(FORM.saveBtn);
  }
  return;
};

export const deleteVersionCmd = function (backend, secretPath, version = 1) {
  return `write ${backend}/delete/${encodePath(secretPath)} versions=${version}`;
};
export const destroyVersionCmd = function (backend, secretPath, version = 1) {
  return `write ${backend}/destroy/${encodePath(secretPath)} versions=${version}`;
};
export const deleteLatestCmd = function (backend, secretPath) {
  return `delete ${backend}/data/${encodePath(secretPath)}`;
};

export const addSecretMetadataCmd = (backend, secret, options = { max_versions: 10 }) => {
  const stringOptions = Object.keys(options).reduce((prev, curr) => {
    return `${prev} ${curr}=${options[curr]}`;
  }, '');
  return `write ${backend}/metadata/${secret} ${stringOptions}`;
};

// Clears kv-related data and capabilities so that admin
// capabilities from setup don't rollover
export function clearRecords(store) {
  store.unloadAll('kv/data');
  store.unloadAll('kv/metatata');
  store.unloadAll('capabilities');
}
