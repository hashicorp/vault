import merge from 'lodash/merge';

// FormField component expects validationMessages in a shape not output by ember-cp-validations
// helper requires validations returned from model validate method
// eg -> const { validations } = await this.model.validate();
export function generateFormFieldErrors(validations) {
  return validations.errors.reduce((errorObj, e) => {
    // check for nested attribute
    if (e.attribute.includes('.')) {
      const keys = e.attribute.split('.').reverse(); // reverse to set value to deepest key first
      const nested = keys.reduce((obj, key) => {
        if (!obj) {
          return { [key]: e.message };
        }
        return { [key]: obj };
      }, null);
      // lodash merge is deep so nested objects will not be replaced
      return merge(errorObj, nested);
    }
    errorObj[e.attribute] = e.message;
    return errorObj;
  }, {});
}
