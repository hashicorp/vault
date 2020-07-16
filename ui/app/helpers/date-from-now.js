import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';

export function dateFromNow([date], options = {}) {
  let d = date;
  try {
    // expects date obj or number only
    return formatDistanceToNow(d, { ...options });
  } catch (e) {
    // if we can't determine the distance, show nothing
    return '';
  }
}

export default helper(dateFromNow);
