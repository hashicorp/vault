import { helper } from '@ember/component/helper';
import formatDate from 'date-fns/format';

export function dateFormat([date, format]) {
  return formatDate(date, format);
}

export default helper(dateFormat);
