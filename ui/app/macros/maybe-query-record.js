/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { computed } from '@ember/object';
import ObjectProxy from '@ember/object/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import { resolve } from 'rsvp';
/**
 * after upgrading to Ember 4.12 a secrets test was erroring with "Cannot create a new tag for `<model::capabilities:undefined>` after it has been destroyed"
 * see this GH issue for information on the fix https://github.com/emberjs/ember.js/issues/16541#issuecomment-382403523
 */
ObjectProxy.reopen({
  unknownProperty(key) {
    if (this.isDestroying || this.isDestroyed) {
      return;
    }

    if (this.content && (this.content.isDestroying || this.content.isDestroyed)) {
      return;
    }

    return this._super(key);
  },
});

export function maybeQueryRecord(modelName, options = {}, ...keys) {
  return computed(...keys, 'store', {
    get() {
      const query = typeof options === 'function' ? options(this) : options;
      const PromiseObject = ObjectProxy.extend(PromiseProxyMixin);

      return PromiseObject.create({
        promise: query ? this.store.queryRecord(modelName, query) : resolve({}),
      });
    },
  });
}
