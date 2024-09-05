/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSubkeysCard
 * @description
sample secret data:
```
 {
  "foo": "abc",
  "bar": {
    "baz": "def"
  },
  "quux": {}
}
```
sample subkeys:
```
 this.subkeys = {
    "bar": {
        "baz": null
    },
    "foo": null,
    "quux": null
}
```
 * 
 * @example
 * <KvSubkeysCard @subkeys={{this.subkeys}} @isPatchAllowed={{true}} />
 *
 * @param {object} subkeys - leaf keys of a kv v2 secret, all values (unless a nested object with more keys) return null
 * @param {boolean} isPatchAllowed - if true, renders the "Patch secret" action. True when: (1) the version is enterprise, (2) a user has "patch" secret + "read" subkeys capabilities, (3) latest secret version is not deleted or destroyed
 */

export default class KvSubkeysCard extends Component {
  @tracked showJson = false;

  @action
  toggleJson(event) {
    this.showJson = event.target.checked;
  }
}
