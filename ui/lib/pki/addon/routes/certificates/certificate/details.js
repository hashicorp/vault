/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiCertificateDetailsRoute extends Route {
  @service api;
  @service secretMountPath;
  @service capabilities;

  async model() {
    const { serial } = this.paramsFor('certificates/certificate');
    const certificate = await this.api.secrets.pkiReadCert(serial, this.secretMountPath.currentPath);
    const { canCreate } = await this.capabilities.for('pkiRevoke', {
      backend: this.secretMountPath.currentPath,
    });

    return {
      certificate: { serial_number: serial, ...certificate },
      canRevoke: canCreate,
    };
  }

  setupController(controller, model) {
    super.setupController(controller, model);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Certificates', route: 'certificates.index', model: this.secretMountPath.currentPath },
      { label: model.certificate.serial_number },
    ];
  }
}
