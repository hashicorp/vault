import { validator, buildValidations } from 'ember-cp-validations';

/**
 * Add validation on dynamic form fields generated via open api spec
 * For fields grouped under default category, add the require/presence validator
 * @param {Array} fieldGroups
 * @returns ember cp validation class
 */
export default function initValidations(fieldGroups) {
  let validators = {};
  fieldGroups.forEach(element => {
    if (element.default) {
      element.default.forEach(v => {
        validators[v.name] = createPresenceValidator(v.name);
      });
    }
  });
  return buildValidations(validators);
}

export const createPresenceValidator = function(label) {
  return validator('presence', {
    presence: true,
    message: `${label} can't be blank.`,
  });
};
