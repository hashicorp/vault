/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type { ExternalOverviewRouteModel } from 'pki/routes/external/overview';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: ExternalOverviewRouteModel;
  isNotConfigured: boolean;
}

export default class ExternalPkiPageOverviewComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked serialNumber = '';
  @tracked orderId = '';

  get countCards() {
    const cards = [];
    const { acmeAccounts, dnsProviders, roles, permissions } = this.args.model;

    const { pkiExternalConfigAcmeAccount, pkiExternalConfigDns, pkiExternalRole } = permissions;
    if (pkiExternalConfigAcmeAccount.canList || pkiExternalConfigAcmeAccount.canRead) {
      cards.push({
        title: 'ACME accounts',
        subtext: 'The total number of ACME accounts (External CAs) configured in this public PKI engine.',
        route: 'external.acme-accounts',
        linkText: 'View ACME accounts',
        count: acmeAccounts.keys.length,
        error: acmeAccounts.errorMsg,
      });
    }

    if (pkiExternalConfigDns.canList || pkiExternalConfigDns.canRead) {
      cards.push({
        title: 'DNS providers',
        subtext: 'The total number of DNS providers configured in this public PKI engine.',
        route: 'external.dns-providers',
        linkText: 'View DNS providers',
        count: dnsProviders.keys.length,
        error: dnsProviders.errorMsg,
      });
    }

    if (pkiExternalRole.canList || pkiExternalRole.canRead) {
      cards.push({
        title: 'Roles',
        subtext: 'The total number of roles configured in this public PKI engine.',
        route: 'external.roles',
        linkText: 'View roles',
        count: roles.keys.length,
        error: roles.errorMsg,
      });
    }

    return cards;
  }

  @action
  lookupCertificate() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.external.certificates.certificate',
      this.args.model.engine.id,
      this.serialNumber
    );
  }

  @action
  lookupOrder() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.external.orders.order',
      this.args.model.engine.id,
      this.orderId
    );
  }
}
