import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

export default AuthConfig.extend({
  useOpenAPI: true,
  organization: attr('string'),
  baseUrl: attr('string', {
    label: 'Base URL',
  }),

  fieldGroups: computed('newFields', function () {
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
