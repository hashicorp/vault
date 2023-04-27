/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

const TOOLS_ACTIONS = ['wrap', 'lookup', 'unwrap', 'rewrap', 'random', 'hash'];

export function toolsActions() {
  return TOOLS_ACTIONS;
}

export default buildHelper(toolsActions);
