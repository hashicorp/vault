/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import KmipRoleForm from 'vault/forms/secrets/kmip/role';

export default class KmipRoleCreateRoute extends Route {
  model() {
    const { scope_name } = this.paramsFor('scope');
    return {
      scopeName: scope_name,
      form: new KmipRoleForm({ operation_all: true }, { isNew: true }),
    };
  }
}
