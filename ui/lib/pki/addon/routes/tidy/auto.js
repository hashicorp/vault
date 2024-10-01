/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default class PkiTidyAutoRoute extends Route {
  model() {
    const { autoTidyConfig } = this.modelFor('tidy');
    return autoTidyConfig;
  }
}
