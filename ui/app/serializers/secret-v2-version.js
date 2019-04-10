import { get } from '@ember/object';
import { assign } from '@ember/polyfills';
import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  secretDataPath: 'data.data',
  normalizeItems(payload) {
    let path = this.secretDataPath;
    // move response that is the contents of the secret from the dataPath
    // to `secret_data` so it will be `secretData` in the model
    payload.secret_data = get(payload, path);
    payload = assign({}, payload, payload.data.metadata);
    delete payload.data;
    payload.path = payload.id;
    // return the payload if it's expecting a single object or wrap
    // it as an array if not
    return payload;
  },
  serialize(snapshot) {
    let secret = snapshot.belongsTo('secret');
    // if both models are stubs, we need to write without CAS
    if (secret.record.isStub && snapshot.record.isStub) {
      return {
        data: snapshot.attr('secretData'),
      };
    }
    let version = secret.record.isStub ? snapshot.attr('version') : secret.attr('currentVersion');
    version = version || 0;
    return {
      data: snapshot.attr('secretData'),
      options: {
        cas: version,
      },
    };
  },
});
