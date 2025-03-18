/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { computed } from '@ember/object';
import Mixin from '@ember/object/mixin';
import { ROUTES } from 'vault/utils/routes';

export default Mixin.create({
  backendCrumb: computed('backend', function () {
    const backend = this.backend;

    if (backend === undefined) {
      throw new Error('backend-crumb mixin requires backend to be set');
    }

    return {
      label: backend,
      text: backend,
      path: ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT,
      model: backend,
    };
  }),
});
