import { helper } from '@ember/component/helper';

export function formatNumber([number]) {
  // formats a number according to the locale
  if (typeof number !== 'number') {
    return;
  }
  return new Intl.NumberFormat().format(number);
}

export default helper(formatNumber);
