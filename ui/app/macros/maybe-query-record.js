/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { computed } from '@ember/object';
import ObjectProxy from '@ember/object/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import { resolve } from 'rsvp';

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
