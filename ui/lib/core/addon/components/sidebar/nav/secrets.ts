/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';
import { RouteName } from 'core/helpers/display-nav-item';

interface Args {
  isEngine?: boolean;
}

export default class SidebarNavSecretsComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;

  routeName = {
    secretsSync: RouteName.SECRETS_SYNC,
  };
}
