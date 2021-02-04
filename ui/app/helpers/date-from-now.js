import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';

export function dateFromNow([date], options = {}) {
  return formatDistanceToNow(date, options);
}

export default helper(dateFromNow);
