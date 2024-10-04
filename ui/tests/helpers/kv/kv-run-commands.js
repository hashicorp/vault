/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-settled-after-test-helper */
import { click, fillIn, visit, settled } from '@ember/test-helpers';
import { FORM } from './kv-selectors';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { assert } from '@ember/debug';
import { kvMetadataPath } from 'vault/utils/kv-path';

// CUSTOM ACTIONS RELEVANT TO KV-V2

export const writeSecret = async function (backend, path, key, val, ns = null) {
  const url = `vault/secrets/${backend}/kv/create`;
  ns ? await visit(url + `?namespace=${ns}`) : await visit(url);
  await settled();
  await fillIn(FORM.inputByAttr('path'), path);
  await fillIn(FORM.keyInput(), key);
  await fillIn(FORM.maskedValueInput(), val);
  await click(FORM.saveBtn);
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
      await fillIn(FORM.keyInput(), key);
      await fillIn(FORM.maskedValueInput(), val);
    } else {
      await fillIn(FORM.keyInput(), `key-${currentVersion}`);
      await fillIn(FORM.maskedValueInput(), `val-${currentVersion}`);
    }
    await click(FORM.saveBtn);
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

// TEST SETUP HELPERS

// sets basic path, backend, and metadata
export const baseSetup = (context) => {
  assert(
    `'baseSetup()' requires mirage: import { setupMirage } from 'ember-cli-mirage/test-support'`,
    context.server
  );
  context.store = context.owner.lookup('service:store');
  context.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
  context.backend = 'kv-engine';
  context.path = 'my-secret';
  context.metadata = metadataModel(context, { withCustom: false });
};

export const metadataModel = (context, { withCustom = false }) => {
  const metadata = withCustom
    ? context.server.create('kv-metadatum', 'withCustomMetadata')
    : context.server.create('kv-metadatum');
  metadata.id = kvMetadataPath(context.backend, context.path);
  context.store.pushPayload('kv/metadata', {
    modelName: 'kv/metadata',
    ...metadata,
  });
  return context.store.peekRecord('kv/metadata', metadata.id);
};
