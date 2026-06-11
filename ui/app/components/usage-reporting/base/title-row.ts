/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type { SafeString } from '@ember/template';

interface TitleRowSignature {
  Args: {
    title: string;
    description?: string | SafeString;
    linkText?: string;
    linkIcon?: string;
    linkUrl?: string;
    linkRoute?: string;
    linkTarget?: '_blank' | '_self';
  };
}

export default class VaultReportingBaseTitleRow extends Component<TitleRowSignature> {
  @service declare readonly router: RouterService;

  get hasExternalLink() {
    return this.args.linkUrl;
  }

  get hasInternalLink() {
    return this.args.linkRoute;
  }

  get linkText() {
    return this.args.linkText || 'View all';
  }

  get linkUrl() {
    return this.args.linkUrl || '#';
  }

  get linkRoute() {
    return this.args.linkRoute;
  }

  get linkIcon() {
    return this.args.linkIcon || 'arrow-right';
  }

  get linkTarget() {
    return this.args.linkTarget || '_self';
  }
}
