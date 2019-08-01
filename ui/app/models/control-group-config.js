import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr } = DS;
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default DS.Model.extend({
  fields: computed(function() {
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
