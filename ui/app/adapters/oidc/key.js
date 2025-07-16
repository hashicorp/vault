/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import NamedPathAdapter from '../named-path';

export default class OidcKeyAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/key';
  }
  rotate(name, verification_ttl) {
    const data = verification_ttl ? { verification_ttl } : {};
    return this.ajax(`${this.urlForUpdateRecord(name, 'oidc/key')}/rotate`, 'POST', { data });
  }
}
