/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import type VersionService from 'vault/services/version';

export default class ClientsController extends Controller {
  @service declare readonly version: VersionService;

  get hasSecretsSync() {
    return this.version.hasSecretsSync;
  }
}
