/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function jsonify([target]) {
  return JSON.parse(target);
}

export default buildHelper(jsonify);
