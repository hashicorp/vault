/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type VersionService from 'vault/services/version';
import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  title: string;
  icon?: string;
  breadcrumbs?: Breadcrumb[];
}

export default class SyncHeaderComponent extends Component<Args> {
  @service declare readonly version: VersionService;

  get badgeTitle() {
    return this.version.isCommunity ? 'Enterprise' : 'Premium Feature';
  }
}
