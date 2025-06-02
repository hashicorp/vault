/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRoute extends Route {
  @service api;

  async model() {
    try {
      const res = await this.api.sys.uiLoginDefaultAuthList(true);
      const loginRules = this.api.keyInfoToArray({ keyInfo: res.keyInfo, keys: res.keys });
      return { loginRules };
    } catch (e) {
      const error = await this.api.parseError(e);
      if (error.status === 404) {
        // If no login settings exist, return an empty array to render the empty state
        return { loginRules: [] };
      }
      // Otherwise fallback to the standard error template
      throw error;
    }
  }
}
