/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiConfigAcmeForm from 'vault/forms/secrets/pki/config/acme';
import PkiConfigClusterForm from 'vault/forms/secrets/pki/config/cluster';
import PkiConfigCrlForm from 'vault/forms/secrets/pki/config/crl';
import PkiConfigUrlsForm from 'vault/forms/secrets/pki/config/urls';

export default class PkiConfigurationEditRoute extends Route {
  @service secretMountPath;

  model() {
    const { acme, cluster, urls, crl, engine, capabilities } = this.modelFor('configuration');
    return {
      engine,
      capabilities,
      acmeForm: new PkiConfigAcmeForm(acme),
      clusterForm: new PkiConfigClusterForm(cluster),
      urlsForm: new PkiConfigUrlsForm(urls),
      crlForm: new PkiConfigCrlForm(crl),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Configuration', route: 'configuration.index', model: this.secretMountPath.currentPath },
      { label: 'Edit' },
    ];
  }
}
