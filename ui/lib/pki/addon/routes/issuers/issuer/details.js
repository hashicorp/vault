/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiIssuerDetailsRoute extends Route {
  @service api;
  @service secretMountPath;
  @service capabilities;

  async fetchCapabilities() {
    const { pathFor } = this.capabilities;
    const backend = this.secretMountPath.currentPath;
    const { issuer_id: issuerId } = this.modelFor('issuers.issuer');

    const pathMap = {
      rotateExported: pathFor('pkiRootRotate', { backend, type: 'exported' }),
      rotateInternal: pathFor('pkiRootRotate', { backend, type: 'internal' }),
      rotateExisting: pathFor('pkiRootRotate', { backend, type: 'existing' }),
      crossSign: pathFor('pkiIntermediateCrossSign', { backend }),
      signIntermediate: pathFor('pkiIssuerSignIntermediate', { backend, issuerId }),
      configure: pathFor('pkiIssuer', { backend, issuerId }),
    };
    const perms = await this.capabilities.fetch(Object.values(pathMap));

    const canRotate =
      perms[pathMap.rotateExported].canUpdate ||
      perms[pathMap.rotateInternal].canUpdate ||
      perms[pathMap.rotateExisting].canUpdate;

    return {
      canRotate,
      canCrossSign: perms[pathMap.crossSign].canUpdate,
      canSignIntermediate: perms[pathMap.signIntermediate].canUpdate,
      canConfigure: perms[pathMap.configure].canUpdate,
    };
  }

  async model() {
    const issuer = this.modelFor('issuers.issuer');
    const { canRotate, canCrossSign, canSignIntermediate, canConfigure } = await this.fetchCapabilities();

    return {
      issuer,
      pem: await this.fetchCertByFormat(issuer.issuer_id, 'pem'),
      der: await this.fetchCertByFormat(issuer.issuer_id, 'der'),
      isRotatable: issuer.isRoot && !!issuer.key_id,
      backend: this.secretMountPath.currentPath,
      canRotate,
      canCrossSign,
      canSignIntermediate,
      canConfigure,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.backend },
      { label: 'Issuers', route: 'issuers.index', model: resolvedModel.backend },
      { label: resolvedModel.issuer.id },
    ];
  }

  /**
   * @private fetches cert by format so it's available for download
   */
  async fetchCertByFormat(issuerId, format) {
    try {
      const path = `/${this.secretMountPath.currentPath}/issuer/${issuerId}/${format}`;
      const response = await this.api.request.get(path);
      const body = format === 'der' ? 'blob' : 'text';
      return response[body]();
    } catch (e) {
      return null;
    }
  }
}
