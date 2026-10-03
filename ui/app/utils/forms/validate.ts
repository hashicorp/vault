/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import validators from 'vault/utils/forms/validators';
import { get } from '@ember/object';

import type { Validations, ValidationMap, FormValidations } from 'vault/app-types';

export const generateErrorCountMessage = (errorCount: number) => {
  if (errorCount < 1) return '';
  // returns count specific message: 'There is an error/are N errors with this form.'
  const isPlural = errorCount > 1 ? `are ${errorCount} errors` : false;
  return `There ${isPlural ? isPlural : 'is an error'} with this form.`;
};

export const validate = (
  data?: unknown,
  validations?: Validations,
  validationMapKey = ''
): FormValidations => {
  let isValid = true;
  const state: ValidationMap = {};
  let errorCount = 0;

  // consider valid when validations are not provided
  if (!validations) {
    return { isValid: true, state, invalidFormMessage: '' };
  }

  for (const key in validations) {
    const rules = validations[key];

    if (!Array.isArray(rules)) {
      console.error(`Must provide validations as an array for property "${key}".`);
      continue;
    }

    // a stateKey may be passed in to map the validations to a nested object
    // eg. validate(this.data, validations, 'data') => { 'data.key': { errors: [], warnings: [], isValid: true } }
    const stateKey = validationMapKey ? `${validationMapKey}.${key}` : key;
    state[stateKey] = { errors: [], warnings: [], isValid: true };

    for (const rule of rules) {
      const { type, options, level, message, validator: customValidator } = rule;
      // check for custom validator or lookup in validators util by type
      const useCustomValidator = typeof customValidator === 'function';
      const validator = useCustomValidator ? customValidator : validators[type];
      if (!validator) {
        console.error(
          !type
            ? 'Validator not found. Either define type or pass custom validator function under "validator" key in validations object'
            : `Validator type: "${type}" not found. Available validators: ${Object.keys(validators).join(
                ', '
              )}`
        );
        continue;
      }
      // dot notation may be used to define key for nested property
      const passedValidation = useCustomValidator
        ? // @ts-expect-error - options may or may not be defined
          validator(data, options)
        : // @ts-expect-error - options may or may not be defined
          validator(get(data, key), options);

      if (!passedValidation) {
        // message can also be a function
        const validationMessage = typeof message === 'function' ? message(data) : message;
        // consider setting a prop like validationErrors directly on the model
        // for now return an errors object
        if (level === 'warn') {
          state[stateKey].warnings.push(validationMessage);
        } else {
          state[stateKey].errors.push(validationMessage);
          if (isValid) {
            isValid = false;
          }
        }
      }
    }
    errorCount += state[stateKey].errors.length;
    state[stateKey].isValid = !state[stateKey].errors.length;
  }

  return { isValid, state, invalidFormMessage: generateErrorCountMessage(errorCount) };
};
