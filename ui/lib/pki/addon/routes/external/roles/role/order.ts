/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';
import timestamp from 'core/utils/timestamp';
import { fetchRoleOrderCert } from 'pki/utils/pki-external-fetch-order';

import type { RoleRouteModel } from '../role';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type RoleOrderRouteModel = ModelFrom<PkiExternalRolesRoleOrderRoute>;

// This route renders the same display component (ExternalPki::OrderCertDetails) as PkiExternalOrdersOrderRoute
// but this order status request is within the context of a role.
// Likely visible to users who only have permissions namespaced to a particular role.
export default class PkiExternalRolesRoleOrderRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async fetchOrderStatus(roleName: string, orderId: string) {
    let details, error;
    try {
      details = await this.api.secrets.pkiExternalCaReadRoleOrderStatus(
        roleName,
        orderId,
        this.secretMountPath.currentPath
      );
    } catch (e) {
      error = await this.api.parseError(e);
      // A 404 here means the user DOES have permission, but the order does not exist.
      // Throw immediately so we don't attempt a cert fetch for a nonexistent order.
      if (error.status === 404) {
        throw error;
      }
    }

    return { details, error };
  }

  async model({ order_id }: { order_id: string }) {
    const { role_name } = this.modelFor('external.roles.role') as RoleRouteModel;
    const order = await this.fetchOrderStatus(role_name, order_id);
    let certificate;
    // Attempt cert fetch if the order completed, or if the order request failed with a non-404
    // since the user may still have permission to fetch the cert directly.
    if (order.details?.order_status === 'completed' || order.error) {
      certificate = await fetchRoleOrderCert(this.api, role_name, order_id, this.secretMountPath.currentPath);
    }

    // If both requests failed, fall back to Ember's standard route error handling.
    if (order?.error && certificate?.error) {
      throw order.error;
    }

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      order_id,
      role_name,
      order,
      certificate,
      responseTimestamp: timestamp.now(),
    };
  }
}
