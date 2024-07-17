/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSubkeys
 * @description
 * ```
 * sample subkeys:
  {
    "bar": {
        "baz": null
    },
    "foo": null,
    "quux": null
}
 * ```
 * 
 * <KvSubkeys @subkeys={{@subkeys}} />
 *
 * @param {object} subkeys - leaf keys of a kv v2 secret, all values (unless a nested object with more keys) return null
 */

export default class KvSubkeys extends Component {
  @tracked showJson = false;
}
