/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiTidyForm from 'vault/forms/secrets/pki/tidy';

export default class PkiTidyManualRoute extends Route {
  @service secretMountPath;

  model() {
    return new PkiTidyForm(
      'PkiTidyRequest',
      {
        acme_account_safety_buffer: 2592000,
        issuer_safety_buffer: 31536000,
        revocation_queue_safety_buffer: 172800,
        safety_buffer: 259200,
        tidy_acme: false,
        tidy_revocation_queue: false,
      },
      { isNew: true }
    );
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Configuration', route: 'configuration.index', model: currentPath },
      { label: 'Tidy', route: 'tidy', model: currentPath },
      { label: 'Manual' },
    ];
  }
}
