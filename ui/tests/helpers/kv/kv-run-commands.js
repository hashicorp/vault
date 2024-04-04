/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-settled-after-test-helper */
import { click, fillIn, visit, settled } from '@ember/test-helpers';
import { KV_FORM } from './kv-selectors';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// CUSTOM ACTIONS RELEVANT TO KV-V2

export const writeSecret = async function (backend, path, key, val, ns = null) {
  const url = `vault/secrets/${backend}/kv/create`;
  ns ? await visit(url + `?namespace=${ns}`) : await visit(url);
  await settled();
  await fillIn(GENERAL.inputByAttr('path'), path);
  await fillIn(KV_FORM.keyInput(), key);
  await fillIn(KV_FORM.maskedValueInput(), val);
  await click(GENERAL.saveButton);
  await settled();
  return;
};

export const writeVersionedSecret = async function (backend, path, key, val, version = 2, ns = null) {
  await writeSecret(backend, path, 'key-1', 'val-1', ns);
  await settled();
  for (let currentVersion = 2; currentVersion <= version; currentVersion++) {
    const url = `/vault/secrets/${backend}/kv/${encodeURIComponent(path)}/details/edit`;
    ns ? await visit(url + `?namespace=${ns}`) : await visit(url);
    await settled();
    if (currentVersion === version) {
      await fillIn(KV_FORM.keyInput(), key);
      await fillIn(KV_FORM.maskedValueInput(), val);
    } else {
      await fillIn(KV_FORM.keyInput(), `key-${currentVersion}`);
      await fillIn(KV_FORM.maskedValueInput(), `val-${currentVersion}`);
    }
    await click(GENERAL.saveButton);
    await settled();
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
