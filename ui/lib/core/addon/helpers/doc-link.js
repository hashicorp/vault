/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { assert } from '@ember/debug';

export function docLink([path]) {
  assert(`doc-link path: "${path}" must begin with a forward slash`, path.startsWith('/'));
  const host = 'https://developer.hashicorp.com';
  return `${host}${path}`;
}

export default helper(docLink);
