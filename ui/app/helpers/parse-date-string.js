import { helper } from '@ember/component/helper';
import { isValid } from 'date-fns';

export function parseDateString(date, separator = '-') {
  // Expects format MM-yyyy by default: no dates
  let datePieces = date.split(separator);
  if (datePieces.length === 2) {
    if (datePieces[0] < 1 || datePieces[0] > 12) {
      throw new Error('Not a valid month value');
    }
    // Since backend converts the timezone to UTC, sending the first (1) as start or end date can cause the month to change.
    // To mitigate this impact of timezone conversion, hard coding the date to avoid month change.
    let date = new Date(datePieces[1], datePieces[0] - 1, 10);
    if (isValid(date)) {
      return date;
    }
  }
  // what to return if not valid?
  throw new Error(`Please use format MM${separator}yyyy`);
}

export default helper(parseDateString);
