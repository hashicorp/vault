/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KubernetesConfigSerializer extends ApplicationSerializer {
  primaryKey = 'backend';

  serialize() {
    const json = super.serialize(...arguments);
    // remove backend value from payload
    delete json.backend;
    // ensure that values from a previous manual configuration are unset
    if (json.disable_local_ca_jwt === false) {
      json.kubernetes_ca_cert = null;
      json.kubernetes_host = null;
      json.service_account_jwt = null;
    }
    return json;
  }
}
