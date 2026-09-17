/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';

export default class UserPreferencesRoute extends Route {
  @service declare readonly flags: FlagsService;

  model() {
    return { isHvdManaged: this.flags.isHvdManaged };
  }
}
