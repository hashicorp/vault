/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import type FlagsService from 'vault/services/flags';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
}
