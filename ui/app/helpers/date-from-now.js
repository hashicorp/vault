import { helper } from '@ember/component/helper';
import { distanceInWordsToNow } from 'date-fns';

export function dateFromNow([date], options = {}) {
  return distanceInWordsToNow(date, options);
}

export default helper(dateFromNow);
