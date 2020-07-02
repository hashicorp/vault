import { helper } from '@ember/component/helper';

export function formatNumber([number]) {
  // formats a number according to the locale
  return new Intl.NumberFormat().format(number);
}

export default helper(formatNumber);
