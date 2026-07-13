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
import type ApiService from 'vault/services/api';
import type { ExternalRouteModel } from '../external';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export interface DnsProviderListItem {
  name: string;
  type: 'aws-route53' | 'azure' | 'google-cloud-dns' | 'rfc2136';
}

export type DnsProvidersRouteModel = ModelFrom<PkiExternalDnsProvidersRoute>;

export default class PkiExternalDnsProvidersRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  supportedDNSTypes = ['aws-route53', 'azure', 'google-cloud-dns', 'rfc2136'] as const;

  fetchConfig(name: string, type: string) {
    const { currentPath } = this.secretMountPath;
    const secretsApi = this.api.secrets;
    switch (type) {
      case 'aws-route53':
        return secretsApi.pkiExternalCaReadConfigDnsAwsRoute53(name, currentPath);
      case 'azure':
        return secretsApi.pkiExternalCaReadConfigDnsAzureDns(name, currentPath);
      case 'google-cloud-dns':
        return secretsApi.pkiExternalCaReadConfigDnsGoogleCloudDns(name, currentPath);
      case 'rfc2136':
        return secretsApi.pkiExternalCaReadConfigDnsRfc2136(name, currentPath);
      default:
        throw Error(
          `Unsupported DNS type. Type must be one of: ${this.supportedDNSTypes.join(', ')}; received: ${type}`
        );
    }
  }

  async model() {
    const { engine, dnsProvidersResp } = this.modelFor('external') as ExternalRouteModel;
    if (dnsProvidersResp.error.message) {
      throw dnsProvidersResp.error;
    }

    // Retrieve type from key_info and request each DNS provider's config
    const infoList = this.api.keyInfoToArray<DnsProviderListItem>(dnsProvidersResp, 'name');
    const results = await Promise.allSettled(
      infoList.map((provider) => this.fetchConfig(provider.name, provider.type))
    );
    const dnsProviders = await Promise.all(
      results.map(async (result, index) => {
        if (result.status === 'fulfilled') {
          return result.value;
        }
        // Edge case: If for some reason user can LIST providers but cannot read
        // config details just return account name
        const { status, message } = await this.api.parseError(result.reason);
        const error =
          status === 403 ? 'You do not have permission to read configurations for this provider' : message;
        return { name: dnsProvidersResp.keys[index], error };
      })
    );

    return { engine, dnsProviders };
  }

  setupController(controller: RouteController, resolvedModel: DnsProvidersRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'DNS providers' },
    ];
  }
}
