import { helper } from '@ember/component/helper';
import { formatDistanceToNow, fromUnixTime } from 'date-fns';

export function dateFromNow([date], options = {}) {
  const time = fromUnixTime(date / 1000);
  return formatDistanceToNow(time, { ...options });
}

export default helper(dateFromNow);
