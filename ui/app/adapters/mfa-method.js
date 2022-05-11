import ApplicationAdapter from './application';

export default class MfaMethodAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    return 'identity/mfa/method';
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
