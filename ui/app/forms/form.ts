/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { validate } from 'vault/utils/forms/validate';
import { set } from '@ember/object';

import type { Validations } from 'vault/app-types';
import type FormField from 'vault/utils/forms/field';

type FormOptions = {
  isNew?: boolean;
};

export default class Form {
  [key: string]: unknown; // Add an index signature to allow dynamic property assignment for set shim
  declare data: Record<string, unknown>;
  declare validations: Validations;
  declare formFields: FormField[];
  declare isNew: boolean;

  constructor(data = {}, options: FormOptions = {}, validations?: Validations) {
    this.data = data;
    this.isNew = options.isNew || false;
    // typically this would be defined on the subclass
    // if validations are conditional, it may be preferable to define them during instantiation
    if (validations) {
      this.validations = validations;
    }
    // to ease migration from Ember Data Models, return a proxy that forwards get/set to the data object for form field props
    // this allows for form field properties to be accessed directly on the class rather than form.data.someField
    return new Proxy(this, {
      get(target, prop: string) {
        const formFields = Array.isArray(target.formFields) ? target.formFields : [];
        const formDataKeys = formFields.map((field) => field.name) || [];
        const getTarget = !Reflect.has(target, prop) && formDataKeys.includes(prop) ? target.data : target;
        return Reflect.get(getTarget, prop);
      },
      set(target, prop: string, value) {
        const formFields = Array.isArray(target.formFields) ? target.formFields : [];
        const formDataKeys = formFields.map((field) => field.name) || [];
        const setTarget = !Reflect.has(target, prop) && formDataKeys.includes(prop) ? target.data : target;
        return Reflect.set(setTarget, prop, value);
      },
    });
  }

  // shim this for now but get away from old Ember patterns!
  set(key = '', val: unknown) {
    set(this, key, val);
  }

  // when overriding in subclass, data can be passed in if serialization of certain values is required
  // this prevents the underlying data object from being mutated causing potential issues in the view
  toJSON(data = this.data) {
    // validate the form
    // if validations are not defined the util will ignore and return valid state
    const formValidation = validate(data, this.validations);
    return { ...formValidation, data };
  }
}
