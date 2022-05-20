import { isPresent } from '@ember/utils';

export const presence = (value) => isPresent(value);

export const length = (value, { nullable = false, min, max } = {}) => {
  if (!min && !max) return;
  // value could be an integer if the attr has a default value of some number
  const valueLength = value?.toString().length;
  if (valueLength) {
    const underMin = min && valueLength < min;
    const overMax = max && valueLength > max;
    return underMin || overMax ? false : true;
  }
  return nullable;
};

export const number = (value, { nullable = false, asString } = {}) => {
  // since 0 is falsy, !value is true even though 0 is valid here
  if (!value && value !== 0) return nullable;
  if (typeof value === 'string' && !asString) {
    return false;
  }
  return !isNaN(value);
};

export default { presence, length, number };
