/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiTidyForm from 'vault/forms/secrets/pki/tidy';

export default class PkiTidyAutoConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const { autoTidyConfig } = this.modelFor('tidy');
    return new PkiTidyForm('PkiConfigureAutoTidyRequest', autoTidyConfig);
  }

  setupController(controller, resolvedModel) {
    // autoTidyConfig id is the backend path
    const { currentPath } = this.secretMountPath;
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Configuration', route: 'configuration.index', model: currentPath },
      { label: 'Tidy', route: 'tidy', model: currentPath },
      { label: 'Auto', route: 'tidy.auto', model: currentPath },
      { label: 'Configure' },
    ];
  }
}
