/* eslint-disable no-console */
import validators from 'vault/utils/validators';
import { get } from '@ember/object';

/**
 * used to validate properties on a class
 *
 * decorator expects validations object with the following shape:
 * { [propertyKeyName]: [{ type, options, message, validator }] }
 * each key in the validations object should refer to the property on the class to apply the validation to
 * type refers to the type of validation to apply -- must be exported from validators util for lookup
 * options is an optional object for given validator -- min, max, nullable etc. -- see validators in util
 * message is added to the errors array and returned from the validate method if validation fails
 * validator may be used in place of type to provide a function that gets executed in the validate method
 * validator is useful when specific validations are needed (dependent on other class properties etc.)
 * validator must be passed as function that takes the class context (this) as the only argument and returns true or false
 * each property supports multiple validations provided as an array -- for example, presence and length for string
 *
 * validations must be invoked using the validate method which is added directly to the decorated class
 * const { isValid, state } = this.model.validate();
 * isValid represents the validity of the full class -- if no properties provided in the validations object are invalid this will be true
 * state represents the error state of the properties defined in the validations object
 * const { isValid, errors } = state[propertyKeyName];
 * isValid represents the validity of the property
 * errors will be populated with messages defined in the validations object when validations fail
 * since a property can have multiple validations, errors is always returned as an array
 *
 *** basic example
 *
 * import Model from '@ember-data/model';
 * import withModelValidations from 'vault/decorators/model-validations';
 *
 * const validations = { foo: [{ type: 'presence', message: 'foo is a required field' }] };
 * @withModelValidations(validations)
 * class SomeModel extends Model { foo = null; }
 *
 * const model = new SomeModel();
 * const { isValid, state } = model.validate();
 * -> isValid = false;
 * -> state.foo.isValid = false;
 * -> state.foo.errors = ['foo is a required field'];
 *
 *** example using custom validator
 *
 * const validations = { foo: [{ validator: (model) => model.bar.includes('test') ? model.foo : false, message: 'foo is required if bar includes test' }] };
 * @withModelValidations(validations)
 * class SomeModel extends Model { foo = false; bar = ['foo', 'baz']; }
 *
 * const model = new SomeModel();
 * const { isValid, state } = model.validate();
 * -> isValid = false;
 * -> state.foo.isValid = false;
 * -> state.foo.errors = ['foo is required if bar includes test'];
 */

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

          state[key] = { errors: [] };

          for (const rule of rules) {
            const { type, options, message, validator: customValidator } = rule;
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
              // consider setting a prop like validationErrors directly on the model
              // for now return an errors object
              state[key].errors.push(message);
              if (isValid) {
                isValid = false;
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
        let isPlural = errorCount > 1 ? `are ${errorCount} errors` : false;
        return `There ${isPlural ? isPlural : 'is an error'} with this form.`;
      }
    };
  };
}
