/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { WIF_ENGINES, allEngines } from 'vault/helpers/mountable-secret-engines';

export default class SecretsBackendConfigurationEditController extends Controller {
  get isWifEngine() {
    return WIF_ENGINES.includes(this.model.type);
  }
  get displayName() {
    return allEngines().find((engine) => engine.type === this.model.type)?.displayName;
  }
}
