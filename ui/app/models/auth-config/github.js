import Ember from 'ember';
import DS from 'ember-data';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
  organization: attr('string'),
  baseUrl: attr('string', {
    label: 'Base URL',
  }),

  fieldGroups: computed(function() {
    const groups = [
      { default: ['organization'] },
      {
        'GitHub Options': ['baseUrl'],
      },
    ];

    return fieldToAttrs(this, groups);
  }),
});
