/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default class BackendIndexRoute extends Route {
  @service router;

  beforeModel() {
    return this.router.replaceWith(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT);
  }
}
