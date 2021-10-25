import { helper } from '@ember/component/helper';
import { format, parseISO } from 'date-fns';

export function dateFormat([date, style], { isFormatted = false, dateOnly = false }) {
  // see format breaking in upgrade to date-fns 2.x https://github.com/date-fns/date-fns/blob/master/CHANGELOG.md#changed-5
  if (isFormatted) {
    return format(new Date(date), style);
  }
  // when date is in '2021-09-01T00:00:00Z' format
  // remove hours so date displays unaffected by timezone
  if (dateOnly) {
    date = date.split('T')[0];
  }
  let number = typeof date === 'string' ? parseISO(date) : date;
  if (!number) {
    return;
  }
  if (number.toString().length === 10) {
    number = new Date(number * 1000);
  }
  return format(number, style);
}

export default helper(dateFormat);
