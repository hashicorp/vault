import DS from 'ember-data';
import Ember from 'ember';
const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const secrets = payload.data.keys.map(secret => {
        let fullSecretPath = payload.id ? payload.id + secret : secret;
        if (!fullSecretPath) {
          fullSecretPath = '\u0020';
        }
        return { id: fullSecretPath };
      });
      return secrets;
    }
    payload.secret_data = payload.data.data;
    delete payload.data;
    return [payload];
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const secrets = nullResponses.includes(requestType) ? { id } : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: secrets };
    // just return the single object because ember is picky
    if (requestType === 'queryRecord') {
      transformedPayload = { [modelName]: secrets[0] };
    }

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serialize(snapshot) {
    return snapshot.attr('secretData');
  },
});
