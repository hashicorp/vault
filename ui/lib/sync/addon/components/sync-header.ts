/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type VersionService from 'vault/services/version';
import type FlagsService from 'vault/services/flags';
import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  title: string;
  icon?: string;
  breadcrumbs?: Breadcrumb[];
}

export default class SyncHeaderComponent extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly flags: FlagsService;

  get badgeText() {
    const isManaged = this.flags.isManaged;
    const onLicense = this.version.hasSecretsSync;
    const isEnterprise = this.version.isEnterprise;

    if (isManaged) return 'Plus feature';
    if (isEnterprise && !onLicense) return 'Premium feature';
    if (!isEnterprise) return 'Enterprise feature';
    // no badge for Enterprise clusters with Secrets Sync on their license--the only remaining option.
    return '';
  }
}
