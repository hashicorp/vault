/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module DashboardWidgetsLearnMore
 * DashboardWidgetsLearnMore component are used to display external links
 *
 * @example
 * ```js
 * <DashboardWidgetsLearnMore  />
 * ```
 */

export default class DashboardWidgetsLearnMore extends Component {
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
