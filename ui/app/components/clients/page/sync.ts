/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import type FlagsService from 'vault/services/flags';
export default class SyncComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;

  title = 'Secrets sync usage';
  description =
    'Secrets sync clients which interacted with Vault for the first time each month. Each bar represents the total new sync clients for that month.';
}
