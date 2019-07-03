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

  nonOperationFields: computed('operationFields', function() {
    let excludeFields = ['role'].concat(this.operationFields);
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
  fieldGroups: computed('fields', 'nonOperationFields', function() {
    const groups = [{ default: this.nonOperationFields }, { 'Allowed Operations': this.operationFields }];
    let ret = fieldToAttrs(this, groups);
    return ret;
  }),

  operationFormFields: computed('operationFieldsWithoutSpecial', function() {
    return expandAttributeMeta(this, this.operationFieldsWithoutSpecial);
  }),
  fields: computed('nonOperationFields', function() {
    return expandAttributeMeta(this, this.nonOperationFields);
  }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`,
});
