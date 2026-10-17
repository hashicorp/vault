/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */
import { validate } from 'vault/utils/forms/validate';
import { set } from '@ember/object';
import { capitalize } from '@ember/string';
import { TrackedObject } from 'tracked-built-ins';
import { maxDurationMessage } from 'core/utils/duration-utils';

import type { Validations } from 'vault/app-types';
import type FormField from 'vault/utils/forms/field';
import type FormFieldGroup from 'vault/utils/forms/field-group';

export type FormOptions = {
  isNew?: boolean;
};

// Flattens the fields declared on a form, including those nested in formFieldGroups.
const collectFormFields = (form: Form<object>): FormField[] => {
  const fields = form.fieldProps.reduce((fields: FormField[], prop) => {
    const formFields = form[prop as keyof typeof form];
    if (Array.isArray(formFields)) {
      fields.push(...formFields);
    }
    return fields;
  }, []);
  const groupFields = form.fieldGroupProps.reduce((groupFields: FormField[], prop) => {
    const formFieldGroups = form[prop as keyof typeof form];
    if (Array.isArray(formFieldGroups)) {
      const fields = formFieldGroups.reduce((arr: FormField[], group: FormFieldGroup) => {
        const values = Object.values(group)[0] || [];
        return [...arr, ...values];
      }, []);
      groupFields.push(...fields);
    }
    return groupFields;
  }, []);
  return [...fields, ...groupFields];
};

// Vault rejects durations it cannot store, so every TTL field gets a max rule on top of the form's own validations.
const withMaxDurationRules = (fields: FormField[], validations?: Validations) => {
  const ttlFields = fields.filter((field) => field.options?.editType === 'ttl');
  if (!ttlFields.length) return validations;

  const merged: Validations = { ...validations };
  // A field can be listed in both formFields and formFieldGroups, so add each rule once per key.
  const seen = new Set<string>();
  for (const { name, options } of ttlFields) {
    const key = options.fieldValue || name;
    if (seen.has(key)) continue;
    seen.add(key);
    // Matches the label FormField renders for unlabeled fields.
    const label = options.label || capitalize(name.replace(/[_-]+/g, ' '));
    merged[key] = [...(merged[key] || []), { type: 'maxDuration', message: maxDurationMessage(label) }];
  }
  return merged;
};

export default class Form<T extends object> {
  declare data: T;
  declare validations: Validations;
  declare isNew: boolean;

  // used by proxy to determine if the property being accessed is a form field
  // override these in subclasses to define additional/different fields defined on the class
  fieldProps = ['formFields'];
  fieldGroupProps = ['formFieldGroups'];

  constructor(data: Partial<T> = {}, options: FormOptions = {}, validations?: Validations) {
    this.data = new TrackedObject(data) as T;
    this.isNew = options.isNew || false;
    // typically this would be defined on the subclass
    // if validations are conditional, it may be preferable to define them during instantiation
    if (validations) {
      this.validations = validations;
    }
    // to ease migration from Ember Data Models, return a proxy that forwards get/set to the data object for form field props
    // this allows for form field properties to be accessed directly on the class rather than form.data.someField
    const proxyTarget = (target: this, prop: string) => {
      try {
        // check if the property that is being accessed is a form field
        const allFields = collectFormFields(target);
        const formDataKeys = allFields.reduce((keys: string[], field) => {
          keys.push(field.name);
          // keyValueFields sub-fields can bind directly to their own model attribute via `valuePath`
          const keyValueFields = field.options?.keyValueFields || [];
          keyValueFields.forEach((keyValueField) => {
            if (keyValueField.valuePath) {
              keys.push(keyValueField.valuePath);
            }
          });
          return keys;
        }, []);
        // if the property is a form field return the data object as the target, otherwise return the original target (this)
        // account for nested form data properties like 'config.maxLeaseTtl' when accessing the object like this.config
        const isDataProp = formDataKeys.some((key) => key === prop || key.split('.').includes(prop));
        return !Reflect.has(target, prop) && isDataProp ? target.data : target;
      } catch (e) {
        // if this fails for any reason return the target object
        return target;
      }
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
    const formValidation = validate(data, withMaxDurationRules(collectFormFields(this), this.validations));
    return { ...formValidation, data };
  }
}
