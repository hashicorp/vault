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
    'This data can be used to understand how many secrets sync clients have been used for this date range. Each Vault secret that is synced to at least one destination counts as one Vault client.';
}
