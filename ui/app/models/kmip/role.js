import DS from 'ember-data';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr, belongsTo } = DS;
const model = DS.Model.extend({
  useOpenAPI: true,
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  updatePath: belongsTo('capabilities'),
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

model.reopenClass({
  relatedCapabilities({ data, included }) {
    let { backend, scope } = data.attributes;
    let url = `${backend}/scope/${scope}/role/${data.id}`;
    data.relationships['updatePath'] = {
      data: {
        type: 'capabilities',
        id: url,
      },
    };

    console.log(data);
    return {
      data,
      included,
    };
  },
});
export default model;
