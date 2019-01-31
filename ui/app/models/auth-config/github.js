import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  organization: attr('string'),
  baseUrl: attr('string', {
    label: 'Base URL',
  }),

  fieldDefinition: computed('newFields', function() {
    const groups = [
      { default: ['organization'] },
      {
        'GitHub Options': ['baseUrl'],
      },
    ];
    if (this.newFields) {
      let allFields = [];
      for (let group in groups) {
        const type = Object.keys(groups[group])[0];
        const field = groups[group][type];
        allFields = allFields.concat(field);
      }
      let otherFields = this.newFields.filter(field => {
        return !allFields.includes(field);
      });
      if (otherFields.length) {
        Object.assign(groups[0].default, groups[0].default.concat(otherFields));
      }
    }

    return groups;
  }),

  fieldGroups: computed('fieldDefinition', function() {
    return this.fieldsToAttrs(this.get('fieldDefinition'));
  }),

  fieldsToAttrs(fieldGroups) {
    return fieldToAttrs(this, fieldGroups);
  },
});
