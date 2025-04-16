/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO pending feedback from the security team, we may keep autocomplete="off" for login fields

import Component from '@glimmer/component';

interface Args {
  loginFields: Field[];
}

interface Field {
  name: string; // sets input name
  label?: string; // label will be "name" capitalized unless label exists
  helperText?: string;
}

export default class AuthFields extends Component<Args> {
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
