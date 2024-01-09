/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable no-console */
import validators from 'vault/utils/validators';
import { get } from '@ember/object';

// see documentation at ui/docs/model-validations.md for detailed usage information
export function withModelValidations(validations) {
  return function decorator(SuperClass) {
    return class ModelValidations extends SuperClass {
      static _validations;

      constructor() {
        super(...arguments);
        if (!validations || typeof validations !== 'object') {
          throw new Error('Validations object must be provided to constructor for setup');
        }
        this._validations = validations;
      }

      validate() {
        let isValid = true;
        const state = {};
        let errorCount = 0;

        for (const key in this._validations) {
          const rules = this._validations[key];

          if (!Array.isArray(rules)) {
            console.error(
              `Must provide validations as an array for property "${key}" on ${this.modelName} model`
            );
            continue;
          }

          state[key] = { errors: [], warnings: [] };

          for (const rule of rules) {
            const { type, options, level, message, validator: customValidator } = rule;
            // check for custom validator or lookup in validators util by type
            const useCustomValidator = typeof customValidator === 'function';
            const validator = useCustomValidator ? customValidator : validators[type];
            if (!validator) {
              console.error(
                !type
                  ? 'Validator not found. Define either type or pass custom validator function under "validator" key in validations object'
                  : `Validator type: "${type}" not found. Available validators: ${Object.keys(
                      validators
                    ).join(', ')}`
              );
              continue;
            }
            const passedValidation = useCustomValidator
              ? validator(this)
              : validator(get(this, key), options); // dot notation may be used to define key for nested property

            if (!passedValidation) {
              // message can also be a function
              const validationMessage = typeof message === 'function' ? message(this) : message;
              // consider setting a prop like validationErrors directly on the model
              // for now return an errors object
              if (level === 'warn') {
                state[key].warnings.push(validationMessage);
              } else {
                state[key].errors.push(validationMessage);
                if (isValid) {
                  isValid = false;
                }
              }
            }
          }
          errorCount += state[key].errors.length;
          state[key].isValid = !state[key].errors.length;
        }

        return { isValid, state, invalidFormMessage: this.generateErrorCountMessage(errorCount) };
      }

      generateErrorCountMessage(errorCount) {
        if (errorCount < 1) return null;
        // returns count specific message: 'There is an error/are N errors with this form.'
        const isPlural = errorCount > 1 ? `are ${errorCount} errors` : false;
        return `There ${isPlural ? isPlural : 'is an error'} with this form.`;
      }
    };
  };
}
