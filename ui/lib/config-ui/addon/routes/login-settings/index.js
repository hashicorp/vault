/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class LoginSettingsRoute extends Route {
  @service api;

  async model() {
    this.createSampleRules();

    const res = await this.api.sys.uiLoginDefaultAuthList(true);
    const loginRules = this.api.keyInfoToArray({ keyInfo: res.keyInfo, keys: res.keys });

    return { loginRules };
  }

  async createSampleRules() {
    await this.api.sys.uiLoginDefaultAuthConfigure('hello', {
      namespace: 'root',
      backup_auth_types: ['okta'],
      default_auth_type: 'ldap',
      disable_inheritance: true,
    });

    await this.api.sys.uiLoginDefaultAuthConfigure('testingTesterz', {
      namespace: 'root',
      backup_auth_types: [],
      default_auth_type: 'oidc',
      disable_inheritance: true,
    });

    await this.api.sys.uiLoginDefaultAuthConfigure('Meowz', {
      namespace: 'root',
      backup_auth_types: ['userpass'],
      default_auth_type: 'userpass',
      DisableInheritance: true,
    });
  }
}
