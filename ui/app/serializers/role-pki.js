import { isNone, isBlank } from '@ember/utils';
import { assign } from '@ember/polyfills';
import { decamelize } from '@ember/string';
import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      let ret = payload.data.keys.map(key => {
        let model = {
          id: key,
        };
        if (payload.backend) {
          model.backend = payload.backend;
        }
        return model;
      });
      return ret;
    }
    assign(payload, payload.data);
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
      isNone(snapshot.changedAttributes()[key])
    ) {
      return;
    }
    if (isBlank(val) && isNone(snapshot.changedAttributes()[key])) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },
  serialize() {
    return this._super(...arguments);
  },
});
