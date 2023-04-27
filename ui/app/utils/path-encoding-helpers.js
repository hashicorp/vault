/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import RouteRecognizer from 'route-recognizer';

const {
  Normalizer: { normalizePath, encodePathSegment },
} = RouteRecognizer;

export function encodePath(path) {
  return path ? path.split('/').map(encodePathSegment).join('/') : path;
}

export { normalizePath, encodePathSegment };
