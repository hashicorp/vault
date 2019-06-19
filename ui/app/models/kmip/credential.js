import DS from 'ember-data';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { computed } from '@ember/object';
const { attr } = DS;

export default DS.Model.extend({
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  role: attr({ readOnly: true }),
  format: attr('string', {
    possibleValues: ['pem', 'der', 'pem_bundle'],
    defaultValue: 'pem',
    label: 'Certificate format',
  }),
  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['format'],
      },
    ];

    return fieldToAttrs(this, groups);
  }),
});
