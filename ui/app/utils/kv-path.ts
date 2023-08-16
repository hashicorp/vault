/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * This set of utils is for calculating the full path for a given KV V2 secret, which doubles as its ID.
 * Additional methods for building URLs for other KV-V2 actions
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';

function buildKvPath(backend: string, path: string, type: string, version?: number | string) {
  const url = `${encodePath(backend)}/${type}/${encodePath(path)}`;
  return version ? `${url}?version=${version}` : url;
}

export function kvDataPath(backend: string, path: string, version?: number | string) {
  return buildKvPath(backend, path, 'data', version);
}
export function kvDeletePath(backend: string, path: string, version?: number | string) {
  return buildKvPath(backend, path, 'delete', version);
}
export function kvMetadataPath(backend: string, path: string) {
  return buildKvPath(backend, path, 'metadata');
}
export function kvDestroyPath(backend: string, path: string) {
  return buildKvPath(backend, path, 'destroy');
}
export function kvUndeletePath(backend: string, path: string) {
  return buildKvPath(backend, path, 'undelete');
}
