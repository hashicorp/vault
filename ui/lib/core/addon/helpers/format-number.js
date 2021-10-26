import { helper } from '@ember/component/helper';

export function formatNumber([number]) {
  if (typeof number !== 'number') {
    return number;
  }
  // formats a number according to the locale
  return new Intl.NumberFormat().format(number);
}

export default helper(formatNumber);
