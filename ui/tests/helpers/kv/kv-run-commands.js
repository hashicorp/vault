/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import listPage from 'vault/tests/pages/secrets/backend/list';

// CUSTOM COMMANDS RELEVANT TO KV-V2

// TODO update writeSecret with create flow from KV ember engine when built
export const writeSecret = async function (backend, path, key, val) {
  await listPage.visitRoot({ backend });
  await listPage.create();
  return editPage.createSecret(path, key, val);
};

export const updateSecret = async function (backend, path, key, val) {
  await editPage.visit({ backend, path });
  await editPage.createNewVersion();
  return editPage.updateSecret(key, val);
};
