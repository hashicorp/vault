import { computed } from '@ember/object';
import ObjectProxy from '@ember/object/proxy';
import PromiseProxyMixin from '@ember/object/promise-proxy-mixin';
import { resolve } from 'rsvp';

export function maybeQueryRecord(modelName, options = {}, ...keys) {
  return computed(...keys, {
    get() {
      const query = typeof options === 'function' ? options(this) : options;
      const PromiseObject = ObjectProxy.extend(PromiseProxyMixin);

      return PromiseObject.create({
        promise: query ? this.get('store').queryRecord(modelName, query) : resolve({}),
      });
    },
  });
}
