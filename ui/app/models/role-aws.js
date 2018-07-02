import Ember from 'ember';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

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
    return expandAttributeMeta(this, keys);
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
