import Ember from 'ember';
import DS from 'ember-data';
import { queryRecord } from 'ember-computed-query';

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

  updatePath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/roles/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  canRead: computed.alias('updatePath.canRead'),

  generatePath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/creds/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canGenerate: computed.alias('generatePath.canUpdate'),

  stsPath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/sts/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canGenerateSTS: computed.alias('stsPath.canUpdate'),
});
