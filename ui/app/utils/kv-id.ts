/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * This util sets the id for the kv/data and kv/metadata data records.
 * @param {string} backend - refers to `backend` attribute on the kv/data model. Examples: kv, secrets, my-kv-engine.
 * @param {number || string} version - refers to `version` attribute on the kv/data model. Examples: '0', 0, 2.
 * @param {string} path - refers to `path` attribute on the kv/data model. Example: my-secret.
 * @param {string} type - either: metadata, data, destroy, or undelete.
 * @returns string id. NOTE: this id is designed to replace the URL for findRecord. Example: my-kv-engine/data/my-secret?=version=2
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';

export function kvId(backend: string, path: string, type: string, version?: string | number) {
  const base = `${encodePath(backend)}/${type}/${encodePath(path)}`;
  return version ? `${base}?version=${version}` : base;
}
