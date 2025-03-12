/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
// import { supportedTypes } from 'vault/helpers/supported-auth-backends';

/**
 * @module Auth::FormTemplate
 *
 * @example
 *
 * @param {string} param - description
 */

export default class AuthFormTemplate extends Component {
  @service version;

  get formComponent() {
    // TODO comment in, mini array below is just for POC
    // const isSupported = supportedTypes(this.version.isEnterprise).includes(this.args.authType);
    const isSupported = ['token', 'okta', 'userpass', 'github'].includes(this.args.authType);
    const component = isSupported ? this.args.authType : 'base';
    // an Auth::Form::<Type> component exists for each type in supported-auth-backends
    // eventually "base" component could be leveraged for rendering custom auth plugins
    return `auth/form/${component}`;
  }
}
