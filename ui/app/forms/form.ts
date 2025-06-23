/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { validate } from 'vault/utils/forms/validate';
import { set } from '@ember/object';

import type { Validations } from 'vault/app-types';
import type FormField from 'vault/utils/forms/field';
import type FormFieldGroup from 'vault/utils/forms/field-group';

export type FormOptions = {
  isNew?: boolean;
};

export default class Form<T extends object> {
  declare data: T;
  declare validations: Validations;
  declare isNew: boolean;

  constructor(data: Partial<T> = {}, options: FormOptions = {}, validations?: Validations) {
    this.data = { ...data } as T;
    this.isNew = options.isNew || false;
    // typically this would be defined on the subclass
    // if validations are conditional, it may be preferable to define them during instantiation
    if (validations) {
      this.validations = validations;
    }
    // to ease migration from Ember Data Models, return a proxy that forwards get/set to the data object for form field props
    // this allows for form field properties to be accessed directly on the class rather than form.data.someField
    const proxyTarget = (target: this, prop: string) => {
      // check if the property that is being accessed is a form field
      const { formFields, formFieldGroups } = target as {
        formFields?: FormField[];
        formFieldGroups?: FormFieldGroup[];
      };
      const fields = Array.isArray(formFields) ? formFields : [];
      // in the case of formFieldGroups we need extract the fields out into a flat array
      const groupFields = Array.isArray(formFieldGroups)
        ? formFieldGroups.reduce((arr: FormField[], group) => {
            const values = Object.values(group)[0] || [];
            return [...arr, ...values];
          }, [])
        : [];
      // combine the formFields and formGroupFields into a single array
      const allFields = [...fields, ...groupFields];
      const formDataKeys = allFields.map((field) => field.name) || [];
      // if the property is a form field return the data object as the target, otherwise return the original target (this)
      // account for nested form data properties like 'config.maxLeaseTtl' when accessing the object like this.config
      const isDataProp = formDataKeys.some((key) => key === prop || key.split('.').includes(prop));

      return !Reflect.has(target, prop) && isDataProp ? target.data : target;
    };

    return new Proxy(this, {
      get(target, prop: string) {
        return Reflect.get(proxyTarget(target, prop), prop);
      },
      set(target, prop: string, value) {
        return Reflect.set(proxyTarget(target, prop), prop, value);
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
