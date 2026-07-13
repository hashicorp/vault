/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb } from 'vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}
export type CertificateRouteModel = ModelFrom<PkiExternalCertificatesCertificateRoute>;

export default class PkiExternalCertificatesCertificateRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model({ serial_number }: { serial_number: string }) {
    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      serial_number,
    };
  }

  setupController(controller: RouteController, resolvedModel: CertificateRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      // There is no "Certificates" index route
      { label: resolvedModel.serial_number },
    ];
  }
}
