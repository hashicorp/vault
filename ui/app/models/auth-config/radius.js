import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
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

  nasIdentifier: attr('string', {
    label: 'NAS Identifier',
  }),

  fieldGroups: computed(function() {
    let groups = [
      {
        default: ['host', 'secret'],
      },
      {
        'RADIUS Options': ['port', 'nasPort', 'nasIdentifier', 'dialTimeout', 'unregisteredUserPolicies'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }

    return fieldToAttrs(this, groups);
  }),
});
