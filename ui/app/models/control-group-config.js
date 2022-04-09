import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  fields: computed(function () {
    return expandAttributeMeta(this, ['maxTtl']);
  }),

  configurePath: lazyCapabilities(apiPath`sys/config/control-group`),
  canDelete: alias('configurePath.canDelete'),
  maxTtl: attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Maximum TTL',
  }),
});
