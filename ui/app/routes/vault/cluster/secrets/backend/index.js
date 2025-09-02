/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class BackendIndexRoute extends Route {
  @service router;

  beforeModel() {
    return this.router.replaceWith('vault.cluster.secrets.backend.list-root');
  }
}
