/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { WIF_ENGINES } from 'vault/utils/all-engines-metadata';

export default class SecretsBackendConfigurationEditController extends Controller {
  get isWifEngine() {
    return WIF_ENGINES.includes(this.model.type);
  }
  get displayName() {
    return engineDisplayData(this.model.type).displayName;
  }
}
