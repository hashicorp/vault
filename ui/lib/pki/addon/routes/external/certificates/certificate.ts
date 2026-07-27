/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';
import timestamp from 'core/utils/timestamp';
import { fetchRoleOrderCert } from 'pki/utils/pki-external-fetch-order';

import type { Breadcrumb } from 'vault/app-types';
import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}
export type CertificateRouteModel = ModelFrom<PkiExternalCertificatesCertificateRoute>;

// Users get here by looking up an order via certificate serial number from the overview certificates card.
// Renders <ExternalPki::OrderCertDetails> (same as PkiExternalOrdersOrderRoute and PkiExternalRolesRoleOrderRoute).
export default class PkiExternalCertificatesCertificateRoute extends Route {
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model({ serial_number }: { serial_number: string }) {
    // This API does not actually return the certificate, just some validity dates and order status
    const certResp = await this.api.secrets.pkiExternalCaReadLookupCert(
      serial_number,
      this.secretMountPath.currentPath
    );

    let certificate;
    const { role_name, order_id, order_status } = certResp;
    // If the order is completed, make request for the actual certificate
    if (role_name && order_id && order_status === 'completed') {
      certificate = await fetchRoleOrderCert(this.api, role_name, order_id, this.secretMountPath.currentPath);
    }

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      serial_number,
      certLookup: certResp,
      certificate,
      responseTimestamp: timestamp.now(),
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
