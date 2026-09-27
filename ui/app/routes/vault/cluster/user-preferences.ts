/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';

export default class UserPreferencesRoute extends Route {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  model() {
    return { isHvdManaged: this.flags.isHvdManaged, isEnterprise: this.version.isEnterprise };
  }
}
