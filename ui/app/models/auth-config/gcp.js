import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  // We have to leave this here because the backend doesn't support the file type yet.
  credentials: attr('string', {
    editType: 'file',
  }),

  googleCertsEndpoint: attr('string'),

  fieldGroups: computed(function() {
    let groups = [
      { default: ['credentials'] },
      {
        'Google Cloud Options': ['googleCertsEndpoint'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    return fieldToAttrs(this, groups);
  }),
});
