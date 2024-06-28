/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import ENV from 'vault/config/environment';

export default class AppFooterComponent extends Component {
  @service version;

  get isDevelopment() {
    return ENV.environment === 'development';
  }
}
