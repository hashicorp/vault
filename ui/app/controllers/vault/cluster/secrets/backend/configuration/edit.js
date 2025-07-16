/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import engineDisplayData from 'vault/helpers/engines-display-data';

export default class SecretsBackendConfigurationEditController extends Controller {
  get isWifEngine() {
    return engineDisplayData(this.model.type)?.isWIF;
  }
  get displayName() {
    return engineDisplayData(this.model.type).displayName;
  }
}
