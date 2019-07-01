import DS from 'ember-data';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import apiPath from 'vault/utils/api-path';
import attachCapabilities from 'vault/lib/attach-capabilities';

const { attr } = DS;
const Model = DS.Model.extend({
  useOpenAPI: true,
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  getHelpUrl(path) {
    return `/v1/${path}/scope/example/role/example?help=1`;
  },

  name: attr({ readOnly: true }),
  fieldGroups: computed(function() {
    let fields = this.newFields.without('role');
    const groups = [{ 'Allowed Operations': fields }];
    return fieldToAttrs(this, groups);
  }),

  fields: computed(function() {
    let fields = this.newFields.removeObjects(['role', 'operationAll', 'operationNone']);
    return expandAttributeMeta(this, fields);
  }),
});

export default attachCapabilities(Model, {
  updatePath: apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`,
});
