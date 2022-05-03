import ApplicationAdapter from './application';

export default class MfaMethodAdapter extends ApplicationAdapter {
  namespace = 'v1';

  urlForQuery(methodType) {
    let baseUrl = this.buildURL() + '/identity/mfa/method';
    if (methodType) {
      return `${baseUrl}/${methodType}`;
    }
    return baseUrl;
  }

  queryRecord(type, id) {
    return this.ajax(this.urlForQuery(type), 'POST', {
      data: {
        id,
      },
    });
  }

  query() {
    return this.ajax(this.urlForQuery(), 'GET', {
      data: {
        list: true,
      },
    }).then((resp) => {
      return resp;
    });
  }
}
