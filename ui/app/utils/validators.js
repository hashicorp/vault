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

export const number = (value, { nullable = false } = {}) => {
  // since 0 is falsy, !value returns true even though 0 is a valid number
  if (!value && value !== 0) return nullable;
  return !isNaN(value);
};

export default { presence, length, number };
