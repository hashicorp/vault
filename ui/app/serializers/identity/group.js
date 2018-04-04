import DS from 'ember-data';
import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend(DS.EmbeddedRecordsMixin, {
  attrs: {
    alias: { embedded: 'always' },
  },

  normalizeFindRecordResponse(store, primaryModelClass, payload) {
    if (payload.alias && Object.keys(payload.alias).length === 0) {
      delete payload.alias;
    }
    return this._super(...arguments);
  },

  serialize() {
    let json = this._super(...arguments);
    delete json.alias;
    if (json.type === 'external') {
      delete json.member_entity_ids;
      delete json.member_group_ids;
    }
    return json;
  },
});
