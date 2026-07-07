/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import {
  SecretsApiPkiExternalCaListConfigAcmeAccountListEnum,
  SecretsApiPkiExternalCaListConfigDnsListEnum,
  SecretsApiPkiExternalCaListRoleListEnum,
} from '@hashicorp/vault-client-typescript';
import { ModelFrom } from 'vault/vault/route';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type { Capabilities } from 'vault/vault/app-types';
import type {
  PkiExternalCaListConfigDnsResponse,
  StandardListResponse,
} from '@hashicorp/vault-client-typescript';

export type ExternalRouteModel = ModelFrom<PkiExternalRoute>;

export default class PkiExternalRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly secretMountPath: SecretMountPath;

  pathKeys = [
    'pkiExternalConfigAcmeAccount',
    'pkiExternalConfigDns',
    'pkiExternalRole',
    'pkiExternalLookupOrders',
  ] as const;

  async fetchPermissions(): Promise<Record<(typeof this.pathKeys)[number], Capabilities>> {
    // Create key/value pair for each key in pathKeys with a value of its API path
    const pathsByKey = this.pathKeys.reduce(
      (obj, key) => {
        obj[key] = this.capabilities.pathFor(key, { backend: this.secretMountPath.currentPath });
        return obj;
      },
      {} as Record<(typeof this.pathKeys)[number], string>
    );

    // Make single API request for capabilities
    const capabilities = await this.capabilities.fetch(Object.values(pathsByKey), {
      routeForCache: 'vault.cluster.secrets.backend.pki.external.overview',
    });

    // Map capabilities back to API path keys so each value is its Capabilities
    return this.pathKeys.reduce(
      (obj, key) => {
        const apiPath = pathsByKey[key];
        const perms = capabilities[apiPath];
        if (perms) obj[key] = perms;
        return obj;
      },
      {} as Record<(typeof this.pathKeys)[number], Capabilities>
    );
  }

  async fetchList(
    listRequest: () => Promise<PkiExternalCaListConfigDnsResponse | StandardListResponse>,
    perms?: Capabilities
  ) {
    let keys: string[] = [],
      errorMsg = '';

    // Only fetch list request if user has permission to avoid additional requests
    if (perms?.canList || perms?.canRead) {
      try {
        const resp = await listRequest();
        keys = resp.keys ?? [];
      } catch (e) {
        // Since this request is only made if a user has permission
        // errors other than a 404 would be unusual, i.e. some sort of internal server issue.
        // Catch just in case to render message in overview card
        const { status, message } = await this.api.parseError(e);
        if (status != 404) {
          errorMsg = message;
        }
      }
    }

    return { keys, errorMsg };
  }

  async model() {
    const perms = await this.fetchPermissions();
    const { currentPath } = this.secretMountPath;
    const [acmeAccounts, dnsProviders, roles] = await Promise.all([
      this.fetchList(
        () =>
          this.api.secrets.pkiExternalCaListConfigAcmeAccount(
            currentPath,
            SecretsApiPkiExternalCaListConfigAcmeAccountListEnum.TRUE
          ),
        perms['pkiExternalConfigAcmeAccount']
      ),
      this.fetchList(
        () =>
          this.api.secrets.pkiExternalCaListConfigDns(
            currentPath,
            SecretsApiPkiExternalCaListConfigDnsListEnum.TRUE
          ),
        perms['pkiExternalConfigDns']
      ),
      this.fetchList(
        () =>
          this.api.secrets.pkiExternalCaListRole(currentPath, SecretsApiPkiExternalCaListRoleListEnum.TRUE),
        perms['pkiExternalRole']
      ),
    ]);

    // If anything at all has been configured, we wan to display that
    const nothingConfigured =
      this.isNotConfigured(perms['pkiExternalConfigAcmeAccount'], acmeAccounts) &&
      this.isNotConfigured(perms['pkiExternalConfigDns'], dnsProviders) &&
      this.isNotConfigured(perms['pkiExternalRole'], roles);

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      acmeAccounts,
      dnsProviders,
      roles,
      permissions: perms,
      isNotConfigured: nothingConfigured,
    };
  }

  isNotConfigured(perms: Capabilities, resp: { keys: string[]; errorMsg: string }) {
    // Only show automation snippets if user has permission to request each resource
    // AND nothing is configured AND there are no errors to display.
    return (perms.canList || perms.canRead) && !resp.keys.length && resp.errorMsg === '';
  }
}
