import { helper } from '@ember/component/helper';
import { formatDistanceToNow, parseISO } from 'date-fns';

export function dateFromNow([date], options = {}) {
  let newDate = typeof date === 'string' ? parseISO(date) : date;
  return formatDistanceToNow(newDate, [options]);
}

export default helper(dateFromNow);
