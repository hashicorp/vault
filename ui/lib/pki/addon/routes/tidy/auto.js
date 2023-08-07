/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

export default class PkiTidyAutoRoute extends Route {
  model() {
    const { autoTidyConfig } = this.modelFor('tidy');
    return autoTidyConfig;
  }
}
