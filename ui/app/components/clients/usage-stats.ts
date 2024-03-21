/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import VersionService from 'vault/services/version';
import { ClientActivityTotals } from 'vault/vault/models/clients/activity';

interface Args {
  totalUsageCounts: ClientActivityTotals;
}

export default class UsageStags extends Component<Args> {
  @service declare readonly version: VersionService;

  get showSecretsSync() {
    return this.version.hasSecretsSync && this.args.totalUsageCounts.secret_syncs > 0;
  }
}
