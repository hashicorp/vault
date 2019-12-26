import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  host: attr('string'),
  secret: attr('string'),

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
