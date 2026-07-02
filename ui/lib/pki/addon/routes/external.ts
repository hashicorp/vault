/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SecretsApiPkiExternalCaListConfigAcmeAccountListEnum } from '@hashicorp/vault-client-typescript';
import { ModelFrom } from 'vault/vault/route';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';

export type ExternalRouteModel = ModelFrom<PkiExternalRoute>;

export default class PkiExternalRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly secretMountPath: SecretMountPath;

  async fetchAcmeAccounts() {
    try {
      const resp = await this.api.secrets.pkiExternalCaListConfigAcmeAccount(
        this.secretMountPath.currentPath,
        SecretsApiPkiExternalCaListConfigAcmeAccountListEnum.TRUE
      );
      return resp.keys;
    } catch {
      // Swallow any other error because we request capabilities separately and
      // this request is just used to determine the mount state for the Overview page.
      return [];
    }
  }

  async model() {
    const { canList } = await this.capabilities.for(
      'pkiExternalCaListConfigAcmeAccount',
      { backend: this.secretMountPath.currentPath },
      { routeForCache: 'vault.cluster.secrets.backend.pki.external' }
    );
    const acmeAccounts = await this.fetchAcmeAccounts();
    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      acmeAccounts,
      canListAcmeConfig: canList,
    };
  }
}
