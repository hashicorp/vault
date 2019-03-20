import { assign } from '@ember/polyfills';
import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeBackend(path, backend) {
    let struct = {};
    for (let attribute in backend) {
      struct[attribute] = backend[attribute];
    }
    //queryRecord adds path to the response
    if (path !== null && !struct.path) {
      struct.path = path;
    }

    if (struct.data) {
      struct = assign({}, struct, struct.data);
      delete struct.data;
    }
    // strip the trailing slash off of the path so we
    // can navigate to it without getting `//` in the url
    struct.id = struct.path.slice(0, -1);
    return struct;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const isCreate = requestType === 'createRecord';
    const isFind = requestType === 'findRecord';
    const isQueryRecord = requestType === 'queryRecord';
    let backends;
    if (isCreate) {
      backends = payload.data;
    } else if (isFind) {
      backends = this.normalizeBackend(id + '/', payload.data);
    } else if (isQueryRecord) {
      backends = this.normalizeBackend(null, payload);
    } else {
      // this is terrible, I'm sorry
      // TODO extract AWS and SSH config saving from the secret-engine model to simplify this
      if (payload.data.secret) {
        backends = Object.keys(payload.data.secret).map(id =>
          this.normalizeBackend(id, payload.data.secret[id])
        );
      } else if (!payload.data.path) {
        backends = Object.keys(payload.data).map(id => this.normalizeBackend(id, payload[id]));
      } else {
        backends = [this.normalizeBackend(payload.data.path, payload.data)];
      }
    }

    return this._super(store, primaryModelClass, backends, id, requestType);
  },

  serialize(snapshot) {
    let type = snapshot.record.get('engineType');
    let data = this._super(...arguments);
    // only KV uses options
    if (type !== 'kv' && type !== 'generic') {
      delete data.options;
    } else if (!data.options.version) {
      // if options.version isn't set for some reason
      // default to 2
      data.options.version = 2;
    }
    return data;
  },
});
