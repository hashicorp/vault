import DS from 'ember-data';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { computed } from '@ember/object';

const { attr } = DS;
export default DS.Model.extend({
  useOpenAPI: true,
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  getHelpUrl(path) {
    return `/v1/${path}/scope/example/role/example?help=1`;
  },

  name: attr('string'),
  allowedOperations: attr(),
  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['name'],
      },
      { 'Allowed Operations': this.newFields },
    ];

    return fieldToAttrs(this, groups);
  }),
});
