/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO pending feedback from the security team, we may keep autocomplete="off" for login fields
// which means deleting this file since `setInputType` can happen directly in the template.

import Component from '@glimmer/component';

/**
 * @module Auth::Fields
 *
 * @example
 * <Auth::Fields @loginFields={{array "username" "password"}} />
 *
 * @param {array} loginFields - array of strings to render as input fields
 */

export default class AuthFields extends Component {
  // token or password should render as "password" types, otherwise render text inputs
  setInputType = (field: string) => (['token', 'password'].includes(field) ? 'password' : 'text');

  setAutocomplete = (field: string) => {
    switch (field) {
      case 'password':
        return 'current-password';
      case 'token':
        return 'off';
      default:
        return field;
    }
  };
}
