/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRoute extends Route {
  @service api;

  async model() {
    const res = await this.api.sys.uiLoginDefaultAuthList(true);
    const loginRules = this.api.keyInfoToArray({ keyInfo: res.keyInfo, keys: res.keys });

    return { loginRules };
  }
}
