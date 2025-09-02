/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function jsonify([target]) {
  // aws secret engine needs to be able to send an empty json value on the field policy_document
  if (!target) return;
  return JSON.parse(target);
}

export default buildHelper(jsonify);
