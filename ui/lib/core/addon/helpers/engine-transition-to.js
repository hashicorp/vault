/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import { service } from '@ember/service';

// this is a helper specifically for routable engines, and assumes that the parent's router is mounted as `app-router`.
export default class EngineTransitionTo extends Helper {
  @service('app-router') router;

  compute(positional, { external = false }) {
    if (external) {
      return () => this.router.transitionToExternal(...positional);
    }
    return () => this.router.transitionTo(...positional);
  }
}
