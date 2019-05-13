import { get } from '@ember/object';
import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  secretDataPath: 'data',
  normalizeItems(payload, requestType) {
    if (requestType !== 'queryRecord' && payload.data.keys && Array.isArray(payload.data.keys)) {
      // if we have data.keys, it's a list of ids, so we map over that
      // and create objects with id's
      return payload.data.keys.map(secret => {
        // secrets don't have an id in the response, so we need to concat the full
        // path of the secret here - the id in the payload is added
        // in the adapter after making the request
        let fullSecretPath = payload.id ? payload.id + secret : secret;

        // if there is no path, it's a "top level" secret, so add
        // a unicode space for the id
        // https://github.com/hashicorp/vault/issues/3348
        if (!fullSecretPath) {
          fullSecretPath = '\u0020';
        }
        return { id: fullSecretPath, backend: payload.backend };
      });
    }
    let path = this.get('secretDataPath');
    // move response that is the contents of the secret from the dataPath
    // to `secret_data` so it will be `secretData` in the model
    payload.secret_data = get(payload, path);
    delete payload[path];
    // return the payload if it's expecting a single object or wrap
    // it as an array if not
    return requestType === 'queryRecord' ? payload : [payload];
  },

  serialize(snapshot) {
    return snapshot.attr('secretData');
  },
});
