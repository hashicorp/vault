import ApplicationAdapter from './application';

export default class MfaMethodAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    return 'identity/mfa/method';
  }

  createOrUpdate(store, type, snapshot) {
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    let id = snapshot.id;
    return this.ajax(`${this.urlForQuery(type.modelName, snapshot)}/${data.type}/${id}`, 'POST', {
      data,
    }).then(() => {
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

  queryRecord(store, type, query) {
    const { id } = query;
    if (!id) {
      throw new Error('MFA method ID is required to fetch the details.');
    }
    const url = this.urlForQuery(query, type.modelName);
    return this.ajax(url, 'POST', {
      data: {
        id,
      },
    });
  }

  query(store, type, query) {
    const url = this.urlForQuery(query, type.modelName);
    return this.ajax(url, 'GET', {
      data: {
        list: true,
      },
    });
  }
}
