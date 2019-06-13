import DS from 'ember-data';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import { computed } from '@ember/object';

const { attr } = DS;
export default DS.Model.extend({
  useOpenAPI: true,
  getHelpUrl(path) {
    return `/v1/${path}/scope/example/role/example?help=1`;
  },

  name: attr('string'),
  allowedOperations: attr(),
  fieldGroups: computed(function() {
    let groups = [
      {
        default: ['name'],
      },
    ];
    groups = combineFieldGroups(groups, this.newFields, []);
    return fieldToAttrs(this, groups);
  }),
});
