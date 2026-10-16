/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import NamespaceForm from 'vault/forms/namespace';

export default class NamespaceCreateRoute extends Route {
  @service version;

  async beforeModel() {
    await this.version.fetchFeatures();
    return super.beforeModel(...arguments);
  }

  model() {
    return this.version.hasNamespaces ? new NamespaceForm({}, { isNew: true }) : null;
  }
}
