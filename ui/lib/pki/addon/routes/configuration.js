/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiConfigurationRoute extends Route {
  @service api;
  @service capabilities;
  @service secretMountPath;

  async fetchCapabilities(backend) {
    const { pathFor } = this.capabilities;
    const pathMap = {
      import: pathFor('pkiIssuersImportBundle', { backend }),
      configAcme: pathFor('pkiConfigAcme', { backend }),
      configCluster: pathFor('pkiConfigCluster', { backend }),
      configCrl: pathFor('pkiConfigCrl', { backend }),
      configUrls: pathFor('pkiConfigUrls', { backend }),
      root: pathFor('pkiRoot', { backend }),
    };
    const perms = await this.capabilities.fetch(Object.values(pathMap));
    return {
      canImportBundle: perms[pathMap.import].canCreate,
      canSetAcme: perms[pathMap.configAcme].canUpdate,
      canSetCluster: perms[pathMap.configCluster].canUpdate,
      canSetCrl: perms[pathMap.configCrl].canUpdate,
      canSetUrls: perms[pathMap.configUrls].canUpdate,
      canDeleteAllIssuers: perms[pathMap.root].canDelete,
    };
  }

  model() {
    const engine = this.modelFor('application');
    const errorHandler = (e) => e.response?.status;
    const { currentPath } = this.secretMountPath;

    return hash({
      engine,
      acme: this.api.secrets
        .pkiReadAcmeConfiguration(currentPath)
        .then((resp) => resp.data) // response type is VoidResponse
        .catch(errorHandler),
      cluster: this.api.secrets.pkiReadClusterConfiguration(currentPath).catch(errorHandler),
      urls: this.api.secrets.pkiReadUrlsConfiguration(currentPath).catch(errorHandler),
      crl: this.api.secrets.pkiReadCrlConfiguration(currentPath).catch(errorHandler),
      capabilities: this.fetchCapabilities(currentPath),
    });
  }
}
