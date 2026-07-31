/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type { Breadcrumb } from 'vault/app-types';
import type { ExternalRouteModel } from '../external';
import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export type AcmeAccountsRouteModel = ModelFrom<PkiExternalAcmeAccountsRoute>;

export default class PkiExternalAcmeAccountsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { engine, acmeAccountsResp } = this.modelFor('external') as ExternalRouteModel;
    if (acmeAccountsResp.error.message) {
      throw acmeAccountsResp.error;
    }

    // Make request for each account in the list
    const results = await Promise.allSettled(
      acmeAccountsResp.keys.map((accountName) =>
        this.api.secrets.pkiExternalCaReadConfigAcmeAccount(accountName, this.secretMountPath.currentPath)
      )
    );
    const acmeAccounts = await Promise.all(
      results.map(async (result, index) => {
        if (result.status === 'fulfilled') {
          return result.value;
        }
        // Edge case: If for some reason user can LIST accounts but cannot read
        // config details just return account name
        const { status, message } = await this.api.parseError(result.reason);
        const error =
          status === 403 ? 'You do not have permission to read configurations for this account' : message;
        return { name: acmeAccountsResp.keys[index], error };
      })
    );

    return { engine, acmeAccounts };
  }

  setupController(controller: RouteController, resolvedModel: AcmeAccountsRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'ACME accounts' },
    ];
  }
}
