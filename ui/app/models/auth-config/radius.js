import Ember from 'ember';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
  host: attr('string'),

  port: attr('number', {
    defaultValue: 1812,
  }),

  secret: attr('string'),

  unregisteredUserPolicies: attr('string', {
    label: 'Policies for unregistered users',
  }),

  dialTimeout: attr('number', {
    defaultValue: 10,
  }),

  nasPort: attr('number', {
    defaultValue: 10,
    label: 'NAS Port',
  }),

  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['host', 'secret'],
      },
      {
        Options: ['port', 'nasPort', 'dialTimeout', 'unregisteredUserPolicies'],
      },
    ];
    return fieldToAttrs(this, groups);
  }),
});
