import DS from 'ember-data';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { computed } from '@ember/object';
const { attr } = DS;
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

const Model = DS.Model.extend({
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  role: attr({ readOnly: true }),
  certificate: attr('string', { readOnly: true }),
  caChain: attr({ readOnly: true }),
  privateKey: attr('string', {
    readOnly: true,
    sensitive: true,
  }),
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

export default attachCapabilities(Model, {
  deletePath: apiPath`${'backend'}/scope/${'scope'}/role/${'role'}/credentials/revoke`,
});
