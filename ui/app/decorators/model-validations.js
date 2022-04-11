/* eslint-disable no-console */
import validators from 'vault/utils/validators';

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

        for (const key in this._validations) {
          const rules = this._validations[key];

          if (!Array.isArray(rules)) {
            console.error(
              `Must provide validations as an array for property "${key}" on ${this.modelName} model`
            );
            continue;
          }

          state[key] = { errors: [] };

          for (const rule of rules) {
            const { type, options, message } = rule;
            if (!validators[type]) {
              console.error(
                `Validator type: "${type}" not found. Available validators: ${Object.keys(validators).join(
                  ', '
                )}`
              );
              continue;
            }
            if (!validators[type](this[key], options)) {
              // consider setting a prop like validationErrors directly on the model
              // for now return an errors object
              state[key].errors.push(message);
              if (isValid) {
                isValid = false;
              }
            }
          }
          state[key].isValid = !state[key].errors.length;
        }
        return { isValid, state };
      }
    };
  };
}
