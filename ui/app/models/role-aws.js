import Ember from 'ember';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr } = DS;
const { computed, get } = Ember;

const CREATE_FIELDS = ['name', 'policy', 'arn'];
export default DS.Model.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  }),
  arn: attr('string', {
    helpText: '',
  }),
  policy: attr('string', {
    helpText: '',
    widget: 'json',
  }),
  attrs: computed(function() {
    let keys = CREATE_FIELDS.slice(0);
    get(this.constructor, 'attributes').forEach((meta, name) => {
      const index = keys.indexOf(name);
      if (index === -1) {
        return;
      }
      keys.replace(index, 1, {
        type: meta.type,
        name,
        options: meta.options,
      });
    });
    return keys;
  }),

  updatePath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  canRead: computed.alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerate: computed.alias('generatePath.canUpdate'),

  stsPath: lazyCapabilities(apiPath`${'backend'}/sts/${'id'}`, 'backend', 'id'),
  canGenerateSTS: computed.alias('stsPath.canUpdate'),
});
