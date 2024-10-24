/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * This set of utils is for calculating the full path for a given KV V2 secret, which doubles as its ID.
 * Additional methods for building URLs for other KV-V2 actions
 */

import { sanitizeStart } from 'core/utils/sanitize-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import queryParamString from './query-param-string';

// only exported for testing
export function buildKvPath(backend: string, path: string, type: string, version?: number | string) {
  const sanitizedPath = sanitizeStart(path); // removing leading slashes
  const url = `${encodePath(backend)}/${type}/${encodePath(sanitizedPath)}`;
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
export function kvSubkeysPath(
  backend: string,
  path: string,
  query: { depth?: number | string; version?: number | string }
) {
  const apiPath = buildKvPath(backend, path, 'subkeys');
  // depth specifies the deepest nesting level the API should return
  // depth=0 returns all subkeys (no limit), depth=1 returns only top-level keys
  const queryParams = queryParamString({
    depth: query?.depth ?? undefined, // no depth returns all levels (no limit)
    version: query?.version ?? undefined, // no version defaults to latest
  });
  return `${apiPath}${queryParams}`;
}
