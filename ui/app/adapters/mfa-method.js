import ApplicationAdapter from './application';

export default class MfaMethodAdapter extends ApplicationAdapter {
  namespace = 'v1';

  methodTypes = ['totp', 'okta', 'duo', 'pingid'];

  pathForType() {
    return 'identity/mfa/method';
  }

  buildURL(methodType) {
    let url = super.buildURL(...arguments);
    return `${url}/${methodType}`;
  }

  async findAll() {
    let results = [];
    for (const type of this.methodTypes) {
      let url = `${this.buildURL(type)}?list=true`;
      try {
        let resp = await this.ajax(url, 'GET');
        results.push(resp.data);
      } catch (err) {
        if (err.httpStatus === 404) {
          // do nothing
        } else {
          throw err;
        }
      }
    }
    return results;
  }

  async findRecord(store, type, id) {
    for (const type of this.methodTypes) {
      let url = `${this.buildURL(type)}/${id}`;
      try {
        let resp = await this.ajax(url, 'GET');
        return resp.data;
      } catch (err) {
        return err;
      }
    }
  }
}
