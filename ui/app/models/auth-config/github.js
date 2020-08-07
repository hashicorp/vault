import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  organization: attr('string'),
  baseUrl: attr('string', {
    label: 'Base URL',
  }),

  fieldGroups: computed(function() {
    let groups = [
      { default: ['organization'] },
      {
        'GitHub Options': ['baseUrl'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }

    return fieldToAttrs(this, groups);
  }),
});
