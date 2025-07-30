/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO this component just exists while the client-list tab is in development for the 1.21 release
// Unless more is added to it, it can be removed when the `client list` tab+route is unhidden
import Component from '@glimmer/component';
import config from 'vault/config/environment';

export default class AuthTabs extends Component {
  get isNotProduction() {
    return config.environment !== 'production';
  }
}
