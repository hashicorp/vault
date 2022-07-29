import ApplicationAdapter from './application';

export default class MfaSetupAdapter extends ApplicationAdapter {
  generate(data) {
    let url = `/v1/identity/mfa/method/totp/generate`;
    return this.ajax(url, 'POST', { data });
  }

  destroy(data) {
    let url = `/v1/identity/mfa/method/totp/destroy`;
    return this.ajax(url, 'POST', { data });
  }
}
