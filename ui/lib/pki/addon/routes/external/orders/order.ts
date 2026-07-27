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
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type OrdersOrderRouteModel = ModelFrom<PkiExternalOrdersOrderRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

// This route renders the same display component (ExternalPki::OrderCertDetails) as PkiExternalRolesRoleOrderRoute
// but uses a global order lookup request instead of one scoped to a role.
// Likely accessible to users with broader, more administrative, permissions.
export default class PkiExternalOrdersOrderRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model({ order_id }: { order_id: string }) {
    // Any error throws and is handled by Ember's route error handling.
    const details = await this.api.secrets.pkiExternalCaReadLookupOrder(
      order_id,
      this.secretMountPath.currentPath
    );

    let certificate;
    // Only fetch the cert if the order has completed and we have a role name to make the request.
    if (details?.role_name && details?.order_status === 'completed') {
      certificate = await fetchRoleOrderCert(
        this.api,
        details.role_name,
        order_id,
        this.secretMountPath.currentPath
      );
    }

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      order_id,
      order: { details },
      certificate,
      responseTimestamp: timestamp.now(),
    };
  }

  setupController(controller: RouteController, resolvedModel: OrdersOrderRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Recent orders', route: 'external.orders', model: currentPath },
      { label: resolvedModel.order_id },
    ];
  }
}
