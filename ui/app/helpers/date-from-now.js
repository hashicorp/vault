import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';
import parseISO from 'date-fns/parseISO';

export function dateFromNow([date], options = {}) {
  let d = date;

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
  }

  try {
    // expects date obj or number only
    return formatDistanceToNow(d, { ...options });
  } catch (e) {
    console.log(e);
    return '';
  }
}

export default helper(dateFromNow);
