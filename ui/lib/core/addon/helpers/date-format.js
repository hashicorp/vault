import { helper } from '@ember/component/helper';
import formatDate from 'date-fns/format';
import parseISO from 'date-fns/parseISO';

export function dateFormat([date, format]) {
  let d = date;

  if (typeof date === 'string') {
    try {
      d = parseISO(date);
    } catch (e) {
      d = new Date(date);
    }
  }

  try {
    console.log('formatting date: ', d);
    console.log('format:', format);
    return formatDate(d, format);
  } catch {
    return 'fooBar';
  }
}

export default helper(dateFormat);
