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

import type { ApiParsedError } from 'vault/vault/api';
import type {
  PkiExternalCaListConfigDnsResponse,
  StandardListResponse,
} from '@hashicorp/vault-client-typescript';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type ExternalRouteModel = ModelFrom<PkiExternalRoute>;

export default class PkiExternalRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async fetchList(listRequest: () => Promise<PkiExternalCaListConfigDnsResponse | StandardListResponse>) {
    let keys: string[] = [],
      key_info = {},
      error: ApiParsedError = { message: '' };

    try {
      const resp = await listRequest();
      keys = resp.keys ?? [];
      if ('key_info' in resp && resp.key_info) {
        key_info = resp.key_info;
      }
    } catch (e) {
      // Catch error to render message in overview card
      // Stored error will be re-thrown in relevant child routes
      error = await this.api.parseError(e);
      if (error.status === 404) {
        // Clear the default message returned by parseError
        // because a 404 is empty and that's expected/okay.
        error.message = '';
      }
    }
    return { keys, key_info, error };
  }

  async model() {
    const { currentPath } = this.secretMountPath;
    const [acmeAccountsResp, dnsProvidersResp, rolesResp] = await Promise.all([
      this.fetchList(() =>
        this.api.secrets.pkiExternalCaListConfigAcmeAccount(
          currentPath,
          SecretsApiPkiExternalCaListConfigAcmeAccountListEnum.TRUE
        )
      ),
      this.fetchList(() =>
        this.api.secrets.pkiExternalCaListConfigDns(
          currentPath,
          SecretsApiPkiExternalCaListConfigDnsListEnum.TRUE
        )
      ),
      this.fetchList(() =>
        this.api.secrets.pkiExternalCaListRole(currentPath, SecretsApiPkiExternalCaListRoleListEnum.TRUE)
      ),
    ]);

    // If a user has permission to request the resource, the request will return either a 404 or data
    // Only show automation snippets if every endpoint returns a 404.
    const showConfigSnippets =
      acmeAccountsResp.error.status === 404 &&
      dnsProvidersResp.error.status === 404 &&
      rolesResp.error.status === 404;

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      acmeAccountsResp,
      dnsProvidersResp,
      rolesResp,
      showConfigSnippets,
    };
  }
}
