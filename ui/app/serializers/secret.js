import Ember from 'ember';
import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  secretDataPath: 'data',
  normalizeItems(payload, requestType) {
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
    let path = this.get('secretDataPath');
    payload.secret_data = Ember.get(payload, path);
    delete payload[path];
    return requestType === 'queryRecord' ? payload : [payload];
  },

  serialize(snapshot) {
    return snapshot.attr('secretData');
  },
});
