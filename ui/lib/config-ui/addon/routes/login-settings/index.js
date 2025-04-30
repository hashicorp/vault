/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default class LoginSettingsRoute extends Route {
  async model() {
    return { loginRules: [{ name: 'Root level auth', namespace: 'root/' }] };
  }
}
