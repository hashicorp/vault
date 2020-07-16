import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';

export function dateFromNow([date], options = {}) {
  // debugger;
  // const time = fromUnixTime(date / 1000);
  let dateString;
  try {
    dateString = formatDistanceToNow(date, { ...options });
  } catch (e) {
    console.log(date);
  }

  return dateString;
}

export default helper(dateFromNow);
