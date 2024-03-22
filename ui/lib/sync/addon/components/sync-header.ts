/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

import type VersionService from 'vault/services/version';
import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  title: string;
  icon?: string;
  breadcrumbs?: Breadcrumb[];
}

export default class SyncHeaderComponent extends Component<Args> {
  @service declare readonly version: VersionService;

  get badgeText() {
    return this.version.hasSecretsSync
      ? ''
      : this.version.isCommunity
      ? 'Enterprise feature'
      : 'Premium feature';
  }
}
