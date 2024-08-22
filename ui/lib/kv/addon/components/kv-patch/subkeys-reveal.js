/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module SubkeysReveal
 *
 * @example
 * <SubkeysReveal @subkeys={{this.subkeys}} />
 *
 * @param {object} subkeys - leaf keys of a kv v2 secret, all values (unless a nested object with more keys) return null
 */

export default class SubkeysReveal extends Component {
  @tracked showSubkeys = false;
}
