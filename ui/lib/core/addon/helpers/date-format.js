import { helper } from '@ember/component/helper';
import formatDate from 'date-fns/format';
import parseISO from 'date-fns/parseISO';

export function dateFormat([date, format]) {
  let d = date;
  let f = format || 'dd MMM yyyy';

  if (typeof date === 'string') {
    // try to parse string assuming ISO format
    d = parseISO(date);

    // if that resulted in invalid date, try with new Date()
    if (!d.getTime()) {
      d = new Date(date);
    }

    // if that also failed, return just passed date string;
    if (!d.getTime()) {
      return date;
    }
  } else if (typeof date === 'object') {
    return '';
  }

  try {
    // expects date obj or number only
    return formatDate(d, f);
  } catch (e) {
    return date;
  }
}

export default helper(dateFormat);
