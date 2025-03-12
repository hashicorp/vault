/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { supportedTypes } from 'vault/helpers/supported-auth-backends';

/**
 * @module Auth::FormTemplate
 *
 * @example
 *
 * @param {string} param - description
 */

export default class AuthFormTemplate extends Component {
  @service version;

  get formFile() {
    const { authType } = this.args;
    if (['oidc', 'jwt'].includes(authType)) return 'oidc-jwt';
    return authType;
  }

  get formComponent() {
    const isSupported = supportedTypes(this.version.isEnterprise).includes(this.args.authType);
    const component = isSupported ? this.formFile : 'base';
    // an Auth::Form::<Type> component exists for each type in supported-auth-backends
    // eventually "base" component could be leveraged for rendering custom auth plugins
    return `auth/form/${component}`;
  }
}
