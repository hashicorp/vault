/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiKeyForm from 'vault/forms/secrets/pki/key';

export default class PkiKeysCreateRoute extends Route {
  @service secretMountPath;

  model() {
    return new PkiKeyForm({}, { isNew: true });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Keys', route: 'keys.index', model: this.secretMountPath.currentPath },
      { label: 'Generate' },
    ];
  }
}
