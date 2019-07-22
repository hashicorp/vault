import DS from 'ember-data';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr } = DS;
export const COMPUTEDS = {
  operationFields: computed('newFields', function() {
    return this.newFields.filter(key => key.startsWith('operation'));
  }),

  operationFieldsWithoutSpecial: computed('operationFields', function() {
    return this.operationFields.slice().removeObjects(['operationAll', 'operationNone']);
  }),

  tlsFields: computed(function() {
    return ['tlsClientKeyBits', 'tlsClientKeyType', 'tlsClientTtl'];
  }),

  nonOperationFields: computed('tlsFields', 'newFields', 'operationFields', function() {
    let excludeFields = ['role'].concat(this.operationFields).concat(this.tlsFields);
    return this.newFields.slice().removeObjects(excludeFields);
  }),
};

const Model = DS.Model.extend(COMPUTEDS, {
  useOpenAPI: true,
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  name: attr({ readOnly: true }),
  getHelpUrl(path) {
    return `/v1/${path}/scope/example/role/example?help=1`;
  },
  fieldGroups: computed('fields', 'tlsFields', 'nonOperationFields', function() {
    const groups = [{ TLS: this.tlsFields }, { 'Allowed Operations': this.operationFields }];
    if (this.nonOperationFields.length) {
      groups.unshift({ default: this.nonOperationFields });
    }
    let ret = fieldToAttrs(this, groups);
    return ret;
  }),

  operationFormFields: computed('operationFieldsWithoutSpecial', function() {
    return expandAttributeMeta(this, this.operationFieldsWithoutSpecial);
  }),
  tlsFormFields: computed('tlsFields', function() {
    return expandAttributeMeta(this, this.tlsFields);
  }),
  fields: computed('nonOperationFields', function() {
    return expandAttributeMeta(this, this.nonOperationFields);
  }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`,
});
