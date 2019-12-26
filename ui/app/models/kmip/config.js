import DS from 'ember-data';
import { computed } from '@ember/object';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { belongsTo } = DS;

export default DS.Model.extend({
  useOpenAPI: true,
  ca: belongsTo('kmip/ca', { async: false }),
  getHelpUrl(path) {
    return `/v1/${path}/config?help=1`;
  },

  fieldGroups: computed(function() {
    let groups = [{ default: ['listenAddrs', 'connectionTimeout'] }];

    groups = combineFieldGroups(groups, this.newFields, []);
    return fieldToAttrs(this, groups);
  }),
});
