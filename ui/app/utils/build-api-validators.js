import { validator, buildValidations } from 'ember-cp-validations';

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
    message: `${label} can't be blank`,
  });
};
