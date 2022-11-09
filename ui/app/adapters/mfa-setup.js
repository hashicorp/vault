import ApplicationAdapter from './application';

export default class MfaSetupAdapter extends ApplicationAdapter {
  adminGenerate(data) {
    const url = `/v1/identity/mfa/method/totp/admin-generate`;
    return this.ajax(url, 'POST', { data });
  }

  adminDestroy(data) {
    const url = `/v1/identity/mfa/method/totp/admin-destroy`;
    return this.ajax(url, 'POST', { data });
  }
}
