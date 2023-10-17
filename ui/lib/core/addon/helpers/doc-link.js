/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function docLink([path]) {
  const host = 'https://developer.hashicorp.com';
  return `${host}${path}`;
}

export default helper(docLink);
