/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type { ApiParsedError } from 'vault/vault/api';
import type { ExternalOverviewRouteModel } from 'pki/routes/external/overview';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: ExternalOverviewRouteModel;
}

export default class ExternalPkiPageOverviewComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked serialNumber = '';
  @tracked orderId = '';

  shouldRenderCard = (resp: { keys: string[]; error: ApiParsedError }) => {
    // Render card if there is data, if the user has permission but the response was empty (404),
    // or if the error message is for something other than 403.
    const status = resp.error['status'];
    return resp.keys.length || status === 404 || (status != 403 && resp.error['message']);
  };

  get countCards() {
    const cards = [];
    const { acmeAccountsResp, dnsProvidersResp, rolesResp } = this.args.model;

    if (this.shouldRenderCard(acmeAccountsResp)) {
      cards.push({
        title: 'ACME accounts',
        subtext: 'The total number of ACME accounts (External CAs) configured in this public PKI engine.',
        route: 'external.acme-accounts',
        linkText: 'View ACME accounts',
        count: acmeAccountsResp.keys.length,
        error: acmeAccountsResp.error.message,
      });
    }

    if (this.shouldRenderCard(dnsProvidersResp)) {
      cards.push({
        title: 'DNS providers',
        subtext: 'The total number of DNS providers configured in this public PKI engine.',
        route: 'external.dns-providers',
        linkText: 'View DNS providers',
        count: dnsProvidersResp.keys.length,
        error: dnsProvidersResp.error.message,
      });
    }

    if (this.shouldRenderCard(rolesResp)) {
      cards.push({
        title: 'Roles',
        subtext: 'The total number of roles configured in this public PKI engine.',
        route: 'external.roles',
        linkText: 'View roles',
        count: rolesResp.keys.length,
        error: rolesResp.error.message,
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
