import { helper } from '@ember/component/helper';
import { differenceInMilliseconds, differenceInSeconds } from 'date-fns';

export default helper(function diffTime([laterDate, earlierDate], { unit = 'milliseconds' }) {
  switch (unit) {
    case 'seconds':
      return differenceInSeconds(laterDate, earlierDate);
    default:
      return differenceInMilliseconds(laterDate, earlierDate);
  }
});
