import { helper } from '@ember/component/helper';
import { isValid } from 'date-fns';

export function parseDateString(date, separator = '-') {
  // Expects format MM-YYYY by default: no dates
  let datePieces = date.split(separator);
  if (datePieces.length === 2) {
    if (datePieces[0] < 1 || datePieces[0] > 12) {
      throw new Error('Not a valid month value');
    }
    let firstOfMonth = new Date(datePieces[1], datePieces[0] - 1, 1);
    if (isValid(firstOfMonth)) {
      return firstOfMonth;
    }
  }
  // what to return if not valid?
  throw new Error(`Please use format MM${separator}YYYY`);
}

export default helper(parseDateString);
