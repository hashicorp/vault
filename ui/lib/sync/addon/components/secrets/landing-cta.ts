/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  isActivated: boolean;
}

export default class LandingCtaComponent extends Component<Args> {
  breadcrumbs = [
    {
      label: 'Vault',
      route: 'vault',
      icon: 'vault',
      linkExternal: true,
    },
    {
      label: 'Secrets sync',
    },
  ];
}
