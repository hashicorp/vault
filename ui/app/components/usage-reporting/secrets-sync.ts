/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toSentenceCase } from 'vault/utils/to-sentence-case';

interface VaultReportingSecretsSyncSignature {
  Args: {
    totalDestinations: number;
    destinations: Record<string, number>;
  };
}

export default class VaultReportingSecretsSync extends Component<VaultReportingSecretsSyncSignature> {
  get hasData() {
    return this.args.totalDestinations > 0;
  }

  get description() {
    if (this.hasData) {
      return 'Total number of destinations (e.g. third-party integrations) synced with secrets';
    }

    return;
  }

  get linkRoute() {
    if (this.hasData) {
      return 'vault.cluster.sync';
    }

    return;
  }

  get totalDestinations() {
    return this.args.totalDestinations || 0;
  }

  get destinationsCountText() {
    const total = this.totalDestinations;
    return `${total} ${total === 1 ? 'destination' : 'destinations'}`;
  }

  get destinationsList() {
    const destinations = this.args.destinations || {};
    return Object.entries(destinations).map(([name, count]) => ({
      name: toSentenceCase(name),
      count,
    }));
  }
}
