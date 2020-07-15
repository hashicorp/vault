import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';
import fromUnixTime from 'date-fns/fromUnixTime';

export function dateFromNow([date], options = {}) {
  const time = fromUnixTime(date);
  // const test = formatDistanceToNow(time);

  return formatDistanceToNow(time, options);
}

export default helper(dateFromNow);
