/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module DashboardLearnMoreCard
 * DashboardLearnMoreCard component are used to display external links
 *
 * @example
 * ```js
 * <DashboardLearnMoreCard  />
 * ```
 */

export default class DashboardLearnMoreCard extends Component {
  get learnMoreLinks() {
    return [
      {
        link: '/vault/tutorials/secrets-management',
        icon: 'docs-link',
        title: 'Secrets Management',
      },
      {
        link: '/vault/tutorials/monitoring',
        icon: 'docs-link',
        title: 'Monitor & Troubleshooting',
      },
      {
        link: '/vault/tutorials/adp/transform',
        icon: 'learn-link',
        title: 'Advanced Data Protection: Transform engine',
        requiredFeature: 'Transform Secrets Engine',
      },
      {
        link: '/vault/tutorials/secrets-management/pki-engine',
        icon: 'learn-link',
        title: 'Build your own Certificate Authority (CA)',
      },
    ];
  }
}
