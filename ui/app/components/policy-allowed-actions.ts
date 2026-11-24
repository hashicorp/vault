/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { AggregatePolicy } from 'vault/utils/policy-aggregator';

interface Args {
  aggregatedPolicy: AggregatePolicy;
  isCard?: boolean;
}

export default class PolicyAllowedActionsComponent extends Component<Args> {
  /**
   * Returns the color class for a capability badge
   * Based on the design:
   * - read, list: gray (neutral)
   * - update, patch: orange/brown (warning)
   * - delete: red (critical)
   * - create, sudo: gray (neutral)
   */
  getCapabilityColor(capability: string): string {
    switch (capability.toLowerCase()) {
      case 'delete':
        return 'critical';
      case 'update':
      case 'patch':
        return 'warning';
      case 'read':
      case 'list':
      case 'create':
      case 'sudo':
      default:
        return 'neutral';
    }
  }

  /**
   * Capitalizes the first letter of a capability for display
   */
  formatCapability(capability: string): string {
    return capability.charAt(0).toUpperCase() + capability.slice(1);
  }

  get sortedPaths(): Array<[string, string[]]> {
    const { policy } = this.args.aggregatedPolicy;
    return Object.entries(policy).sort(([pathA], [pathB]) => pathA.localeCompare(pathB));
  }
}
