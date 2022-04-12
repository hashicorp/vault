import { isPresent } from '@ember/utils';

export const presence = (value) => isPresent(value);

export const length = (value, { nullable = false, min, max } = {}) => {
  let isValid = nullable;
  if (typeof value === 'string') {
    const underMin = min && value.length < min;
    const overMax = max && value.length > max;
    isValid = underMin || overMax ? false : true;
  }
  return isValid;
};

export const number = (value, { nullable = false, asString } = {}) => {
  if (!value) return nullable;
  if (typeof value === 'string' && !asString) {
    return false;
  }
  return !isNaN(value);
};

export default { presence, length, number };
