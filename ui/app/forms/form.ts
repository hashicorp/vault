/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { validate } from 'vault/utils/forms/validate';
import { set } from '@ember/object';

import type { Validations } from 'vault/app-types';

export default class Form {
  declare data;
  declare validations: Validations;

  constructor(data = {}, validations?: Validations) {
    this.data = data;
    // typically this would be defined on the subclass
    // if validations are conditional, it may be preferable to define them during instantiation
    if (validations) {
      this.validations = validations;
    }
  }

  // shim this for now but get away from old Ember patterns!
  set(key: string, val: unknown) {
    set(this, key, val);
  }

  // when overriding in subclass, data can be passed in if serialization of certain values is required
  // this prevents the underlying data object of being mutated causing potential issues in the view
  toJSON(data = this.data) {
    // validate the form
    // if validations are not defined the util will ignore and return valid state
    const formValidation = validate(data, this.validations, 'data');
    return { ...formValidation, data };
  }
}
