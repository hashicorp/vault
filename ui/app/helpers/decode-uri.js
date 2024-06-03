/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function decodeUri(string) {
  return decodeURI(string);
}

export default buildHelper(decodeUri);
