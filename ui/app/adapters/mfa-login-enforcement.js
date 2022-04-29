import ApplicationAdapter from '../application';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';

  pathForType() {
    return 'identity/mfa/login-enforcement';
  }

  async query(store, type, query) {
    const url = this.buildURL(type.modelName, null, null, 'query', query);
    return this.ajax(url, 'GET', { data: { list: true } }).then((resp) => resp.data);
  }
}
