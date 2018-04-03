import DS from 'ember-data';
import Ember from 'ember';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import AuthConfig from '../../auth-config';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
  safetyBuffer: attr({
    defaultValue: '72h',
    editType: 'ttl',
  }),

  disablePeriodicTidy: attr('boolean', {
    defaultValue: false,
  }),

  attrs: computed(function() {
    return expandAttributeMeta(this, ['safetyBuffer', 'disablePeriodicTidy']);
  }),
});
