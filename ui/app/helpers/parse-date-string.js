import { helper } from '@ember/component/helper';
import { isValid } from 'date-fns';

export function parseDateString([date], separator = '-') {
  // Expects format MM-YYYY by default: no dates
  let datePieces = date.split(separator);
  if (datePieces.length > 1) {
    let startDate = new Date(Date.UTC(datePieces[1], datePieces[0] - 1, 1));
    if (isValid(startDate)) {
      return startDate;
    }
  }
  // what to return if not valid?
  return null;
}

export default helper(parseDateString);
