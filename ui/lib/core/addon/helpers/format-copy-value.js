/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

// Hds::CopyButton value must be a string to be copied
// template helper takes any value and returns a copyable string
export function formatCopyValue([value]) {
  if (!value) return value;
  const type = Array.isArray(value) ? 'array' : typeof value;
  switch (type) {
    case 'string':
      return value;
    case 'object':
      return JSON.stringify(value);
    case 'array':
      return value.join('\n');
    default:
      return value.toString();
  }
}

export default buildHelper(formatCopyValue);
