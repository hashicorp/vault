/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface Args {
  backend: SecretsEngineResource;
  showConfigSnippets: boolean;
}

export default class ExternalPkiHeaderTabsComponent extends Component<Args> {
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
}
