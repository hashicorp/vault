import DS from 'ember-data';
import Ember from 'ember';
const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  extractLazyPaginatedData(payload) {
    let ret;
    if (payload.zero_address_roles) {
      payload.zero_address_roles.forEach(role => {
        // mutate key_info object to add zero_address info
        payload.data.key_info[role].zero_address = true;
      });
    }
    ret = payload.data.keys.map(key => {
      let model = {
        id: key,
        key_type: payload.data.key_info[key].key_type,
        zero_address: payload.data.key_info[key].zero_address,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
    delete payload.data.key_info;
    return ret;
  },

  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys;
    }
    Ember.assign(payload, payload.data);
    delete payload.data;
    return [payload];
  },
  modelNameFromPayloadKey(payloadType) {
    return payloadType;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const responseJSON = nullResponses.includes(requestType) ? { id } : this.normalizeItems(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: responseJSON };
    // just return the single object because ember is picky
    if (requestType === 'queryRecord') {
      transformedPayload = { [modelName]: responseJSON[0] };
    }

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    if (attributes.options.readOnly) {
      return;
    }
    if (
      attributes.type === 'object' &&
      val &&
      Object.keys(val).length > 0 &&
      Ember.isNone(snapshot.changedAttributes()[key])
    ) {
      return;
    }
    if (Ember.isBlank(val) && Ember.isNone(snapshot.changedAttributes()[key])) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },
  serialize() {
    return this._super(...arguments);
  },
});
