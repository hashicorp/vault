/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import { FORM } from './kv-selectors';

// CUSTOM COMMANDS RELEVANT TO KV-V2

export const writeSecret = async function (backend, path, key, val) {
  await visit(`vault/secrets/${backend}/kv/create`);
  await fillIn(FORM.inputByAttr('path'), path);
  await fillIn(FORM.keyInput(), key);
  await fillIn(FORM.maskedValueInput(), val);
  return click(FORM.saveBtn);
};

export const writeVersionedSecret = async function (backend, path, key, val, version = 2) {
  await writeSecret(backend, path, 'key-1', 'val-1');
  for (let currentVersion = 2; currentVersion <= version; currentVersion++) {
    await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent(path)}/details/edit`);
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

export const addSecretMetadataCmd = (backend, secret, options = { max_versions: 10 }) => {
  const stringOptions = Object.keys(options).reduce((prev, curr) => {
    return `${prev} ${curr}=${options[curr]}`;
  }, '');
  return `write ${backend}/metadata/${secret} ${stringOptions}`;
};
