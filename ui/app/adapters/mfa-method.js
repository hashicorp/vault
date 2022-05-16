import ApplicationAdapter from './application';

export default class MfaMethodAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    return 'identity/mfa/method';
  }

  createOrUpdate(store, type, snapshot) {
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    const { id } = snapshot;
    return this.ajax(this.buildURL(type.modelName, data.type, id, 'POST'), 'POST', {
      data,
    }).then(() => {
      // TODO: Check how 204's are handled by ember
      return {
        data: {
          id,
          ...data,
        },
      };
    });
  }

  createRecord() {
    return this.createOrUpdate(...arguments);
  }

  updateRecord() {
    return this.createOrUpdate(...arguments);
  }

  query(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    return this.ajax(url, 'GET', {
      data: {
        list: true,
      },
    });
  }

  buildURL(modelName, type, id, requestType) {
    if (requestType === 'POST') {
      return `${super.buildURL(modelName, null, null, null)}/${type}/${id}`;
    }
    return super.buildURL(...arguments);
  }
}
