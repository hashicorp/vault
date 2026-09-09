/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { matchesCurrentUrl } from 'core/helpers/matches-current-url';

import type { Breadcrumb } from 'vault/vault/app-types';
import type { PATH_MAP } from 'vault/utils/constants/capabilities';
import type CapabilitiesService from 'vault/services/capabilities';
import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface Args {
  backend: SecretsEngineResource;
  showConfigSnippets: boolean;
}

type PathKeys = (keyof typeof PATH_MAP)[];

// Maps external PKI sub-routes (everything after "vault.cluster.secrets.backend.pki.")
// to the API path keys relevant for the policy flyout.
const EXTERNAL_ROUTE_PATH_MAP: Record<string, PathKeys> = {
  'external.overview': [
    'pkiExternalConfigAcmeAccount',
    'pkiExternalConfigDns',
    'pkiExternalRoleList',
    'pkiExternalLookupOrders',
  ],
  'external.orders.index': ['pkiExternalLookupOrdersRecent'],
  'external.orders.order': ['pkiExternalLookupOrder', 'pkiExternalRoleOrderFetchCert'],
  'external.certificates.certificate': ['pkiExternalLookupCert', 'pkiExternalRoleOrderFetchCert'],
  'external.dns-providers': ['pkiExternalConfigDns'],
  'external.acme-accounts': ['pkiExternalConfigAcmeAccount'],
  'external.roles.index': ['pkiExternalRoleList'],
  'external.roles.role.overview': [
    'pkiExternalRole',
    'pkiExternalRoleActiveOrders',
    'pkiExternalRoleCachedCert',
  ],
  'external.roles.role.order': ['pkiExternalRoleOrderStatus', 'pkiExternalRoleOrderFetchCert'],
  'external.roles.role.active-orders': ['pkiExternalRoleActiveOrders'],
};

interface Args {
  title: string;
  breadcrumbs: Breadcrumb[];
  backend: SecretsEngineResource;
  roleName: string;
}

export default class ExternalPkiHeaderTabsComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly capabilities: CapabilitiesService;

  get policyPaths() {
    const backend = this.args.backend?.id;
    const keys = this._routePathKeys(this.router.currentRouteName);

    const roleName = this.args.roleName || ':role_name';
    return keys ? this.capabilities.pathsForList(keys, { backend, roleName }) : null;
  }

  get defaultTabs() {
    return this.args.showConfigSnippets
      ? [{ label: 'Overview', route: 'external.overview' }]
      : [
          { label: 'Overview', route: 'external.overview' },
          { label: 'Roles', route: 'external.roles' },
          { label: 'Recent orders', route: 'external.orders' },
          { label: 'DNS providers', route: 'external.dns-providers' },
          { label: 'ACME accounts', route: 'external.acme-accounts' },
        ];
  }

  private _routePathKeys(routeName: string | null | undefined): PathKeys | undefined {
    const prefix = 'vault.cluster.secrets.backend.pki.';
    let suffix = routeName?.startsWith(prefix) ? routeName?.slice(prefix.length) : routeName;
    if (suffix === 'external.error' || suffix === 'external.roles.role.error') {
      // Find the matching suffix via URL instead
      suffix = Object.keys(EXTERNAL_ROUTE_PATH_MAP).find((k) => matchesCurrentUrl(this.router, prefix + k));
    }
    return suffix ? EXTERNAL_ROUTE_PATH_MAP[suffix] : undefined;
  }
}
