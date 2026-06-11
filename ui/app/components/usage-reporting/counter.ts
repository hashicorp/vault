/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';

interface VaultReportingCounterSignature {
  Args: {
    count: number;
    title: string;
    tooltipMessage?: string;
    icon?: string;
    suffix?: string;
    link?: string;
    emptyText?: string;
    emptyLink?: string;
  };
}

export default class VaultReportingCounter extends Component<VaultReportingCounterSignature> {
  @service declare readonly router: RouterService;

  get shouldShowEmptyState() {
    return this.args.count === 0 && this.args.emptyText;
  }

  get count() {
    if (this.shouldShowEmptyState) {
      return this.args.emptyText;
    }

    if (this.args.suffix) {
      return `${this.args.count} ${this.args.suffix}`;
    }

    return this.args.count;
  }

  get link() {
    const routeToCheck = this.shouldShowEmptyState ? this.args.emptyLink : this.args.link;
    return routeToCheck;
  }
}
