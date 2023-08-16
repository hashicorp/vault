/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default class MfaSetupAdapter extends ApplicationAdapter {
  currentTokenGenerate(data) {
    const url = `/v1/identity/mfa/method/totp/generate`;
    return this.ajax(url, 'POST', { data });
  }

  adminDestroy(data) {
    const url = `/v1/identity/mfa/method/totp/admin-destroy`;
    return this.ajax(url, 'POST', { data });
  }
}
